package errs

import "net/http"

var (
	ErrInvalidCredentials = New(http.StatusUnauthorized, "invalid email or password")
	ErrEmailExists        = New(http.StatusBadRequest, "email already registered")
	ErrInvalidToken       = New(http.StatusUnauthorized, "invalid refresh token")
	ErrExpiredToken       = New(http.StatusUnauthorized, "expired refresh token")
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
