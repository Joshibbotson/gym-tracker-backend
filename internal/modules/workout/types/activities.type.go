package types

import (
	"time"
)

type DailyWorkout struct {
	Date          time.Time      `json:"date"`
	WorkoutConfig *WorkoutConfig `json:"workoutConfig,omitempty"`
}

type MonthlyData struct {
	Month    string         `json:"month"`
	Workouts []DailyWorkout `json:"workouts"`
}

type YearlyData struct {
	Year   int           `json:"year"`
	Months []MonthlyData `json:"months"`
}

type WorkoutData struct {
	Data []YearlyData `json:"data"`
}
