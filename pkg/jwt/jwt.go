package jwt

import (
	"time"

	"Postpartum_BackEnd/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func getSecret() []byte {
	return []byte(utils.GetEnv("JWT_SECRET"))
}

func GenerateAccessToken(userID uuid.UUID, role, name, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"name":    name,
		"email":   email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}
