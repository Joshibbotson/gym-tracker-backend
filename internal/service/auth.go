package service

import (
	"context"
	"time"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// handle business logic for auth
// register, login, validateToken,

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// AuthService defines methods for user authentication actions
type AuthService interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(name, email, password string) (*User, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

// (r *authService) this is a method receiver it's like a class and this is it's method
func (r *authService) CreateUser(name string, email string, password string) (*User, error) {
	collection := db.Client.Database("gym-tracker").Collection("user")

	user := User{
		Name:     name,
		Email:    email,
		Password: password, //hash this really
	}

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (r *authService) GetUserByEmail(email string) (*User, error) {
	collection := db.Client.Database("gym-tracker").Collection("user")

	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
