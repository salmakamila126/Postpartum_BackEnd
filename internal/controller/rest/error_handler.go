package rest

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/errs"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, dto.Error(appErr.Message))
		return
	}

	c.JSON(http.StatusInternalServerError, dto.Error("internal server error"))
}
