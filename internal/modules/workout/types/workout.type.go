package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	ID        primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Date      time.Time          `bson:"date" json:"date"`
	Workout   *WorkoutConfig     `bson:"workout,omitempty" json:"workout"` // Pointer to make it nullable
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}
