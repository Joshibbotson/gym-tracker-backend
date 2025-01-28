package user

import (
	"github.com/joshibbotson/gym-tracker-backend/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
}

type userRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository() UserRepository {
	return &userRepository{
		userCollection: db.Client.Database(db.DB_NAME).Collection("user"),
	}
}

func FetchUserDetails()
