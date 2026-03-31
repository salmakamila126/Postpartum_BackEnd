package usecase

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/entity"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/repository"
	"net/http"

	"github.com/google/uuid"
)

type UserUsecase struct {
	Repo *repository.Repository
}

func NewUserUsecase(repo *repository.Repository) *UserUsecase {
	return &UserUsecase{Repo: repo}
}

func (u *UserUsecase) GetProfile(userID uuid.UUID) (*entity.User, error) {
	return u.Repo.UserRepository.FindByID(userID)
}

func (u *UserUsecase) UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) (*entity.User, error) {
	user, err := u.Repo.UserRepository.FindByID(userID)
	if err != nil {
		return nil, errs.New(http.StatusBadRequest, "user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := u.Repo.UserRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
