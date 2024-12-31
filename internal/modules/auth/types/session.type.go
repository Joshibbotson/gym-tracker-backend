package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	SessionID string             `bson:"session_id"`
	ExpiresAt time.Time          `bson:"expires_at"`
}
