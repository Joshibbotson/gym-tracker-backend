package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/auth/types"

	"github.com/google/uuid"
	db "github.com/joshibbotson/gym-tracker-backend/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GetUserByEmail(email string) (*t.User, error)
	CreateUser(name, email, password string) (*t.User, error)
	Login(email, password string) (*t.Session, error)
	createOrUpdateSession(userID primitive.ObjectID, name string, email string) (t.Session, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

// (r *authService) this is a method receiver it's like a class and this is it's method
func (r *authService) CreateUser(name string, email string, password string) (*t.User, error) {
	collection := db.Client.Database(db.DB_NAME).Collection("user")

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
	user := t.User{
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

// should return a cookie perhaps instead of User?
func (r *authService) Login(email string, password string) (*t.Session, error) {
	collection := db.Client.Database(db.DB_NAME).Collection("user")

	// Set a timeout for the database query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the user by email
	var user t.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No user found for email:", email) // Debugging
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	session, err := r.createOrUpdateSession(user.ID, user.Name, user.Email)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (r *authService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *authService) GetUserByEmail(email string) (*t.User, error) {
	collection := db.Client.Database(db.DB_NAME).Collection("user")

	var user t.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (*authService) createOrUpdateSession(userID primitive.ObjectID, name string, email string) (t.Session, error) {
	sessionCollection := db.Client.Database(db.DB_NAME).Collection("session")
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	session := t.Session{
		UserID:    userID,
		Name:      name,
		Email:     email,
		SessionID: sessionID,
		ExpiresAt: expiresAt,
	}

	update := bson.M{
		"$set": bson.M{
			"expires_at": expiresAt,
		},
	}

	err := sessionCollection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"user_id": session.UserID},
		update).Decode(&session)
	if err != nil {
		fmt.Printf("FindOneAndUpdate error: %v\n", err)
	} else {
		fmt.Printf("Existing session: %+v\n", session)
	}

	if err == nil {
		return session, nil
	}

	// if no session available insert one.
	_, err = sessionCollection.InsertOne(context.TODO(), session)
	if err != nil {
		return t.Session{}, err
	}

	return session, nil
}
