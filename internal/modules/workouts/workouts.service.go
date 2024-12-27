package workouts

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutType string

const (
	WorkoutTypePush      WorkoutType = "push"
	WorkoutTypePull      WorkoutType = "pull"
	WorkoutTypeLegs      WorkoutType = "legs"
	WorkoutTypeArms      WorkoutType = "arms"
	WorkoutTypeChest     WorkoutType = "chest"
	WorkoutTypeShoulders WorkoutType = "shoulders"
	WorkoutTypeBack      WorkoutType = "back"
)

type CaloriePhase string

const (
	CaloriePhaseCut      CaloriePhase = "cut"
	CaloriePhaseBulk     CaloriePhase = "bulk"
	CaloriePhaseMaintain CaloriePhase = "maintain"
)

type WorkoutConfig struct {
	Weight       *float64      `json:"weight,omitempty" bson:"weight,omitempty"`
	WorkoutType  *WorkoutType  `json:"workoutType,omitempty" bson:"workoutType,omitempty"`
	CaloriePhase *CaloriePhase `json:"caloriePhase,omitempty" bson:"caloriePhase,omitempty"`
	ChestSize    *float64      `json:"chestSize,omitempty" bson:"chestSize,omitempty"`
	WaistSize    *float64      `json:"waistSize,omitempty" bson:"waistSize,omitempty"`
	BicepSize    *float64      `json:"bicepSize,omitempty" bson:"bicepSize,omitempty"`
	ForearmSize  *float64      `json:"forearmSize,omitempty" bson:"forearmSize,omitempty"`
	ThighSize    *float64      `json:"thighSize,omitempty" bson:"thighSize,omitempty"`
	CalfSize     *float64      `json:"calfSize,omitempty" bson:"calfSize,omitempty"`
}

type Workout struct {
	ID   primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Date time.Time          `bson:"date" json:"date"`

	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt"`
}

type CreateWorkout struct {
}

// handle business logic for workouts
type WorkoutsService interface {
	createWorkout(workout CreateWorkout) (*Workout, error)
}
