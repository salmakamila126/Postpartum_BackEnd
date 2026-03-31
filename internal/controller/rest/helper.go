package rest

import (
	"Postpartum_BackEnd/internal/errs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	raw, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errs.New(http.StatusUnauthorized, "unauthorized")
	}

	str, ok := raw.(string)
	if !ok {
		return uuid.Nil, errs.New(http.StatusUnauthorized, "invalid user_id in token")
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, errs.New(http.StatusUnauthorized, "invalid user_id format")
	}

	return id, nil
}
