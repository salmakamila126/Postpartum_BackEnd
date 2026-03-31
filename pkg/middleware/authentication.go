package middleware

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/pkg/logger"
	"Postpartum_BackEnd/pkg/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func getSecret() []byte {
	return []byte(utils.GetEnv("JWT_SECRET"))
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.Error("authorization header missing"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, dto.Error("invalid authorization header format"))
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return getSecret(), nil
		})

		if err != nil || !token.Valid {
			logger.Log.Warn("invalid token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, dto.Error("invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.Error("invalid token claims"))
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		if name, ok := claims["name"].(string); ok {
			c.Set("user_name", name)
		}
		if email, ok := claims["email"].(string); ok {
			c.Set("user_email", email)
		}

		c.Next()
	}
}
