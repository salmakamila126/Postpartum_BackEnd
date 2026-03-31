package errs

import "net/http"

var (
	ErrPsychologistNotFound = New(http.StatusNotFound, "psychologist not found")
)
