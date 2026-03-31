package rest

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PsychologistController struct {
	Usecase *usecase.PsychologistUsecase
}

func NewPsychologistController(u *usecase.PsychologistUsecase) *PsychologistController {
	return &PsychologistController{Usecase: u}
}

func (c *PsychologistController) GetAll(ctx *gin.Context) {
	list, err := c.Usecase.GetAll()
	if err != nil {
		HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Success(list))
}

func (c *PsychologistController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		HandleError(ctx, errs.New(http.StatusBadRequest, "invalid psychologist id"))
		return
	}

	detail, err := c.Usecase.GetByID(id)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Success(detail))
}

func (c *PsychologistController) UpdatePhotoURL(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		HandleError(ctx, errs.New(http.StatusBadRequest, "invalid psychologist id"))
		return
	}

	var req dto.UpdatePhotoURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, errs.New(http.StatusBadRequest, "photo_url is required"))
		return
	}

	if err := c.Usecase.UpdatePhotoURL(id, req.PhotoURL); err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Success("photo updated"))
}

func (c *PsychologistController) BookingWhatsApp(ctx *gin.Context) {
	psychologistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		HandleError(ctx, errs.New(http.StatusBadRequest, "invalid psychologist id"))
		return
	}

	if _, err := GetUserID(ctx); err != nil {
		HandleError(ctx, err)
		return
	}

	userNameRaw, _ := ctx.Get("user_name")
	userEmailRaw, _ := ctx.Get("user_email")
	userName, _ := userNameRaw.(string)
	userEmail, _ := userEmailRaw.(string)

	var req dto.BookingWhatsAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, errs.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	result, err := c.Usecase.BuildBookingWhatsApp(psychologistID, userName, userEmail, req)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Success(result))
}
