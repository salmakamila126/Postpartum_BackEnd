package config

type SleepConfig struct {
	MaxSessionsPerDay int
	MinDaysForPredict int
	MaxRecentData     int
	PredictCycles     int
}

func NewSleepConfig() SleepConfig {
	return SleepConfig{
		MaxSessionsPerDay: 8,
		MinDaysForPredict: 3,
		MaxRecentData:     50,
		PredictCycles:     3,
	}
}
