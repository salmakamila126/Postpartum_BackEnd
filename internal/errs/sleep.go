package errs

import "net/http"

var (
	ErrNoData                 = New(http.StatusBadRequest, "no data")
	ErrNotEnoughData          = New(http.StatusBadRequest, "not enough data")
	ErrLastSessionNotFinished = New(http.StatusBadRequest, "last session not finished")
	ErrAlreadySleeping        = New(http.StatusBadRequest, "baby already sleeping")
	ErrNoActiveSession        = New(http.StatusBadRequest, "no active sleep session")

	ErrBackdateModify = New(http.StatusBadRequest, "cannot modify backdate data")
	ErrMaxSession     = New(http.StatusBadRequest, "maximum 8 sleep sessions per day reached")
	ErrOverlap        = New(http.StatusBadRequest, "sleep session overlap")
	ErrInvalidTime    = New(http.StatusBadRequest, "end must be after start")

	ErrManualOnlyToday  = New(http.StatusBadRequest, "manual input only allowed for today")
	ErrBulkOnlyPast     = New(http.StatusBadRequest, "bulk input only allowed for past dates")
	ErrBulkAlreadyExist = New(http.StatusBadRequest, "data already exists for this date")

	ErrPredictionFailed = New(http.StatusInternalServerError, "prediction failed")
)
