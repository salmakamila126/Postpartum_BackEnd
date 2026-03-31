package dto

import "time"

type SleepPredictionItem struct {
	Sleep string `json:"sleep"`
	Wake  string `json:"wake"`
}

type SleepPredictionResponse struct {
	Predictions []SleepPredictionItem `json:"predictions"`
}

type SleepSessionItem struct {
	ID              string `json:"id"`
	Start           string `json:"start"`
	End             string `json:"end,omitempty"`
	DurationMinutes int    `json:"duration_minutes"`
	IsBackdate      bool   `json:"is_backdate"`
}

type DailySleepResponse struct {
	Date              string             `json:"date"`
	TotalSleepMinutes int                `json:"total_sleep_minutes"`
	TotalSessions     int                `json:"total_sessions"`
	AvgSleepMinutes   int                `json:"avg_sleep_minutes"`
	Sessions          []SleepSessionItem `json:"sessions"`
	IsSleeping        bool               `json:"is_sleeping"`
	CurrentStart      string             `json:"current_start,omitempty"`
}

type InsightResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SleepManualRequest struct {
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

type SleepBulkRequest struct {
	Date     string               `json:"date" binding:"required"`
	Sessions []SleepManualRequest `json:"sessions" binding:"required,min=1"`
}

type SleepStatusResponse struct {
	Sleeping bool   `json:"sleeping"`
	Start    string `json:"start,omitempty"`
}

type SleepHistoryItem struct {
	ID              string    `json:"id"`
	SleepTime       string    `json:"sleep_time"`
	WakeTime        *string   `json:"wake_time"`
	DurationMinutes *int      `json:"duration_minutes"`
	IsBackdate      bool      `json:"is_backdate"`
	CreatedAt       time.Time `json:"created_at"`
}

type SleepHistoryResponse struct {
	Sessions []SleepHistoryItem `json:"sessions"`
}
