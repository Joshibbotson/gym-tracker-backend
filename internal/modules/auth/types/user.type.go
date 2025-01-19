package types

import (
	"time"

	a "github.com/joshibbotson/gym-tracker-backend/internal/modules/auth/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password,omitempty" json:"password,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`

	// AuthData fields (optional)
	Surname       string          `bson:"surname,omitempty" json:"surname,omitempty"`
	FirstName     string          `bson:"firstName,omitempty" json:"firstName,omitempty"`
	AuthId        string          `bson:"authId,omitempty" json:"authId,omitempty"`
	PictureUrl    string          `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty"`
	VerifiedEmail bool            `bson:"verifiedEmail,omitempty" json:"verifiedEmail,omitempty"`
	AuthProvider  a.AuthProviders `bson:"authProvider,omitempty" json:"authProvider,omitempty"`
}
