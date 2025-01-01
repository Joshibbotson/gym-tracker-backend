package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Date      time.Time          `bson:"date" json:"date"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
	Workout   *WorkoutConfig     `bson:"workout,omitempty" json:"workout"` // Pointer to make it nullable
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
