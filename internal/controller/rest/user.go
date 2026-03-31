package rest

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Usecase *usecase.UserUsecase
}

func NewUserController(u *usecase.UserUsecase) *UserController {
	return &UserController{Usecase: u}
}

func (uc *UserController) Profile(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	user, err := uc.Usecase.GetProfile(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		BirthDate: user.Baby.BirthDate,
	}))
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	user, err := uc.Usecase.UpdateProfile(userID, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		BirthDate: user.Baby.BirthDate,
	}))
}
