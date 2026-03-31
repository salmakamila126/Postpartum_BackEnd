package rest

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SymptomController struct {
	Usecase *usecase.SymptomUsecase
}

func NewSymptomController(u *usecase.SymptomUsecase) *SymptomController {
	return &SymptomController{Usecase: u}
}

func (c *SymptomController) CreateOrUpdate(ctx *gin.Context) {
	var req dto.CreateSymptomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, errs.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	userID, err := GetUserID(ctx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	result, err := c.Usecase.CreateOrUpdate(userID, req)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Success(dto.SymptomSaveResponse{
		Date:  req.Date,
		Alert: usecase.ToResponse(result),
	}))
}

func (c *SymptomController) GetHistory(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	data, err := c.Usecase.GetHistory(userID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Success(data))
}

func (c *SymptomController) GetDetail(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	dateStr := ctx.Param("date")
	if dateStr == "" {
		HandleError(ctx, errs.New(http.StatusBadRequest, "date param is required"))
		return
	}

	data, err := c.Usecase.GetDetail(userID, dateStr)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.Success(data))
}
