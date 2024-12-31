package types

import (
	c "github.com/joshibbotson/gym-tracker-backend/internal/modules/workout/constants"
)

type WorkoutConfig struct {
	Weight       *float64        `json:"weight,omitempty" bson:"weight,omitempty"`
	WorkoutType  *c.WorkoutType  `json:"workoutType,omitempty" bson:"workoutType,omitempty"`
	CaloriePhase *c.CaloriePhase `json:"caloriePhase,omitempty" bson:"caloriePhase,omitempty"`
	ChestSize    *float64        `json:"chestSize,omitempty" bson:"chestSize,omitempty"`
	WaistSize    *float64        `json:"waistSize,omitempty" bson:"waistSize,omitempty"`
	BicepSize    *float64        `json:"bicepSize,omitempty" bson:"bicepSize,omitempty"`
	ForearmSize  *float64        `json:"forearmSize,omitempty" bson:"forearmSize,omitempty"`
	ThighSize    *float64        `json:"thighSize,omitempty" bson:"thighSize,omitempty"`
	CalfSize     *float64        `json:"calfSize,omitempty" bson:"calfSize,omitempty"`
}
