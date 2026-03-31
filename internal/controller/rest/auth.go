package rest

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Usecase *usecase.AuthUsecase
}

func NewAuthController(u *usecase.AuthUsecase) *AuthController {
	return &AuthController{Usecase: u}
}

func (ac *AuthController) Register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	user, accessToken, refreshToken, err := ac.Usecase.Register(input)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.Success(dto.AuthResponse{
		User: dto.UserResponse{
			Name:      user.Name,
			Email:     user.Email,
			BirthDate: user.Baby.BirthDate,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}))
}

func (ac *AuthController) Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	user, access, refresh, err := ac.Usecase.Login(input)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.AuthResponse{
		User: dto.UserResponse{
			Name:      user.Name,
			Email:     user.Email,
			BirthDate: user.Baby.BirthDate,
		},
		AccessToken:  access,
		RefreshToken: refresh,
	}))
}

func (ac *AuthController) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "refresh_token is required"))
		return
	}

	access, refresh, err := ac.Usecase.RefreshToken(body.RefreshToken)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}))
}

func (ac *AuthController) Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "refresh_token is required"))
		return
	}

	if err := ac.Usecase.Logout(body.RefreshToken); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.Success("logged out successfully"))
}
