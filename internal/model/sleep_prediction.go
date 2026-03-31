package model

import "time"

type SleepPrediction struct {
	NextSleep time.Time `json:"next_sleep"`
	NextWake  time.Time `json:"next_wake"`
}
