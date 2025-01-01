package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyWorkout struct {
	ID     primitive.ObjectID `json:"_id"`
	Date   time.Time          `json:"date"`
	Config *WorkoutConfig     `json:"workoutConfig,omitempty"`
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
