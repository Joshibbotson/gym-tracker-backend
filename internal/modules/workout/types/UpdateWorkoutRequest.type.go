package types

import (
	"time"

	c "github.com/joshibbotson/gym-tracker-backend/internal/modules/workout/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateWorkoutRequest struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	UserId           primitive.ObjectID `bson:"userID" json:"userID"`
	Date             time.Time          `bson:"date" json:"date"`
	Weight           *float64           `json:"weight,omitempty" bson:"weight,omitempty"`
	TargetMuscles    []c.TargetMuscles  `json:"targetMuscles,omitempty" bson:"targetMuscles,omitempty"`
	CaloriePhase     *c.CaloriePhase    `json:"caloriePhase,omitempty" bson:"caloriePhase,omitempty"`
	NeckSize         *float64           `json:"neckSize,omitempty" bson:"neckSize,omitempty"`
	ShoulderSize     *float64           `json:"shoulderSize,omitempty" bson:"shoulderSize,omitempty"`
	LeftCalfSize     *float64           `json:"leftCalfSize,omitempty" bson:"leftCalfSize,omitempty"`
	RightCalfSize    *float64           `json:"rightCalfSize,omitempty" bson:"rightCalfSize,omitempty"`
	LeftAnkleSize    *float64           `json:"leftAnkleSize,omitempty" bson:"leftAnkleSize,omitempty"`
	RightAnkleSize   *float64           `json:"rightAnkleSize,omitempty" bson:"rightAnkleSize,omitempty"`
	LeftThighSize    *float64           `json:"leftThighSize,omitempty" bson:"leftThighSize,omitempty"`
	RightThighSize   *float64           `json:"rightThighSize,omitempty" bson:"rightThighSize,omitempty"`
	LeftWristSize    *float64           `json:"leftWristSize,omitempty" bson:"leftWristSize,omitempty"`
	RightWristSize   *float64           `json:"rightWristSize,omitempty" bson:"rightWristSize,omitempty"`
	ChestSize        *float64           `json:"chestSize,omitempty" bson:"chestSize,omitempty"`
	WaistSize        *float64           `json:"waistSize,omitempty" bson:"waistSize,omitempty"`
	HipSize          *float64           `json:"hipSize,omitempty" bson:"hipSize,omitempty"`
	LeftBicepSize    *float64           `json:"leftBicepSize,omitempty" bson:"leftBicepSize,omitempty"`
	RightBicepSize   *float64           `json:"rightBicepSize,omitempty" bson:"rightBicepSize,omitempty"`
	LeftForearmSize  *float64           `json:"leftForearmSize,omitempty" bson:"leftForearmSize,omitempty"`
	RightForearmSize *float64           `json:"rightForearmSize,omitempty" bson:"rightForearmSize,omitempty"`
}
