package usecase

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/repository"
	"Postpartum_BackEnd/internal/usecase/service"
	"Postpartum_BackEnd/pkg/jwt"
	"Postpartum_BackEnd/pkg/logger"
	"Postpartum_BackEnd/pkg/timeutil"
	"Postpartum_BackEnd/pkg/utils"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	Repo *repository.Repository
	DB   *gorm.DB
}

func NewAuthUsecase(repo *repository.Repository, db *gorm.DB) *AuthUsecase {
	return &AuthUsecase{Repo: repo, DB: db}
}

func (u *AuthUsecase) Register(input dto.RegisterRequest) (*entity.User, string, string, error) {
	if input.Password != input.ConfirmPassword {
		return nil, "", "", errs.New(http.StatusBadRequest, "password and confirm_password do not match")
	}
	if _, err := timeutil.ParseDate(input.BirthDate); err != nil {
		return nil, "", "", errs.New(http.StatusBadRequest, "invalid birth_date format - use YYYY-MM-DD")
	}

	var user entity.User
	var refreshToken string

	err := u.DB.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewRepository(tx)

		existingUser, err := txRepo.UserRepository.FindByEmail(input.Email)
		if err == nil && existingUser != nil {
			return errs.ErrEmailExists
		}

		hashed, err := service.HashPassword(input.Password)
		if err != nil {
			return err
		}

		user = entity.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: hashed,
		}

		if err := txRepo.UserRepository.Create(&user); err != nil {
			return err
		}

		baby := entity.Baby{
			UserID:    user.UserID,
			BirthDate: input.BirthDate,
		}

		if err := txRepo.BabyRepository.Create(&baby); err != nil {
			return err
		}

		_, refreshToken, err = service.GenerateTokens(user.UserID, user.Role, user.Name, user.Email)
		if err != nil {
			return err
		}

		hashedToken := utils.HashToken(refreshToken)

		rt := entity.RefreshToken{
			Token:     hashedToken,
			UserID:    user.UserID,
			ExpiresAt: timeutil.NowWIB().Add(7 * 24 * time.Hour),
		}

		return txRepo.RefreshTokenRepository.Create(&rt)
	})

	if err != nil {
		logger.Log.Error("register transaction failed", zap.String("email", input.Email), zap.Error(err))
		return nil, "", "", err
	}

	accessToken, err := jwt.GenerateAccessToken(user.UserID, user.Role, user.Name, user.Email)
	if err != nil {
		logger.Log.Error("generate access token failed", zap.String("user_id", user.UserID.String()), zap.Error(err))
		return nil, "", "", err
	}

	return &user, accessToken, refreshToken, nil
}

func (u *AuthUsecase) Login(input dto.LoginRequest) (*entity.User, string, string, error) {
	user, err := u.Repo.UserRepository.FindByEmail(input.Email)
	if err != nil {
		logger.Log.Warn("login failed — email not found", zap.String("email", input.Email))
		return nil, "", "", errs.ErrInvalidCredentials
	}

	if err := service.ComparePassword(user.Password, input.Password); err != nil {
		logger.Log.Warn("login failed — wrong password", zap.String("email", input.Email))
		return nil, "", "", errs.ErrInvalidCredentials
	}

	access, refresh, err := service.GenerateTokens(user.UserID, user.Role, user.Name, user.Email)
	if err != nil {
		logger.Log.Error("generate token failed", zap.String("user_id", user.UserID.String()), zap.Error(err))
		return nil, "", "", err
	}

	hashed := utils.HashToken(refresh)
	rt := entity.RefreshToken{
		Token:     hashed,
		UserID:    user.UserID,
		ExpiresAt: timeutil.NowWIB().Add(7 * 24 * time.Hour),
	}

	if err := u.Repo.RefreshTokenRepository.Create(&rt); err != nil {
		logger.Log.Error("save refresh token failed", zap.String("user_id", user.UserID.String()), zap.Error(err))
		return nil, "", "", err
	}

	return user, access, refresh, nil
}

func (u *AuthUsecase) RefreshToken(token string) (string, string, error) {
	hashed := utils.HashToken(token)

	rt, err := u.Repo.RefreshTokenRepository.Find(hashed)
	if err != nil {
		logger.Log.Error("refresh token not found", zap.Error(err))
		return "", "", errs.ErrInvalidToken
	}

	if timeutil.NowWIB().After(rt.ExpiresAt) {
		logger.Log.Warn("refresh token expired", zap.String("user_id", rt.UserID.String()))
		if err := u.Repo.RefreshTokenRepository.Delete(hashed); err != nil {
			logger.Log.Error("failed to delete expired refresh token", zap.Error(err))
		}
		return "", "", errs.ErrExpiredToken
	}

	user, err := u.Repo.UserRepository.FindByID(rt.UserID)
	if err != nil {
		logger.Log.Error("user not found during token refresh", zap.String("user_id", rt.UserID.String()), zap.Error(err))
		return "", "", err
	}

	if err := u.Repo.RefreshTokenRepository.Delete(hashed); err != nil {
		logger.Log.Error("failed to delete old refresh token", zap.Error(err))
	}

	access, newRefresh, err := service.GenerateTokens(user.UserID, user.Role, user.Name, user.Email)
	if err != nil {
		return "", "", err
	}

	newHashed := utils.HashToken(newRefresh)
	newRT := entity.RefreshToken{
		Token:     newHashed,
		UserID:    user.UserID,
		ExpiresAt: timeutil.NowWIB().Add(7 * 24 * time.Hour),
	}

	if err := u.Repo.RefreshTokenRepository.Create(&newRT); err != nil {
		return "", "", err
	}

	return access, newRefresh, nil
}

func (u *AuthUsecase) Logout(token string) error {
	hashed := utils.HashToken(token)
	return u.Repo.RefreshTokenRepository.Delete(hashed)
}
