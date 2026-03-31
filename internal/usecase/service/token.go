package service

import (
	"Postpartum_BackEnd/pkg/jwt"

	"github.com/google/uuid"
)

func GenerateTokens(userID uuid.UUID, role, name, email string) (string, string, error) {
	access, err := jwt.GenerateAccessToken(userID, role, name, email)
	if err != nil {
		return "", "", err
	}

	refresh := uuid.New().String()
	return access, refresh, nil
}
