package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// handle business logic for auth
// register, login, validateToken,

const DB_NAME = "gym-tracker"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// AuthService defines methods for user authentication actions
type AuthService interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(name, email, password string) (*User, error)
	Login(email, password string) (*User, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

// (r *authService) this is a method receiver it's like a class and this is it's method
func (r *authService) CreateUser(name string, email string, password string) (*User, error) {
	collection := db.Client.Database(DB_NAME).Collection("user")

	// Check if a user with the email already exists
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := r.HashPassword(password)
	if err != nil {
		return nil, err
	}
	fmt.Println("Generated hash during user creation:", hashedPassword)
	user := User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (r *authService) Login(email string, password string) (*User, error) {
	collection := db.Client.Database(DB_NAME).Collection("user")

	// Set a timeout for the database query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Debugging: Log email to ensure correct email is received
	fmt.Println("Attempting login for email:", email)

	// Find the user by email
	var user User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No user found for email:", email) // Debugging
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Debugging: Log found user’s email
	fmt.Println("User found in database with email:", user.Email)

	fmt.Println("Stored hashed password:", user.Password)
	fmt.Println("Entered password:", password)

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Password comparison error:", err) // Debugging
		return nil, errors.New("incorrect password")
	}

	// Authentication successful
	fmt.Println("User authenticated successfully:", email) // Debugging
	return &user, nil
}

func (r *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (r *authService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *authService) GetUserByEmail(email string) (*User, error) {
	collection := db.Client.Database(DB_NAME).Collection("user")

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

func (r *authService) GenerateToken() (string, error) {

}
