package errs

import "net/http"

var (
	ErrSymptomNotFound = New(http.StatusNotFound, "symptom record not found")
	ErrOnlyTodayUpdate = New(http.StatusForbidden, "past records cannot be edited")
	ErrBackdateExists  = New(http.StatusConflict, "a record already exists for this date; past records are immutable")
	ErrFutureDate      = New(http.StatusBadRequest, "cannot submit records for a future date")
	ErrInvalidUser     = New(http.StatusUnauthorized, "unauthorized")
)
