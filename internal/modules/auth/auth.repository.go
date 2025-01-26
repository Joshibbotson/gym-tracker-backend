package auth

import (
	"context"
	"time"

	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/auth/types"

	db "github.com/joshibbotson/gym-tracker-backend/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	FindUserByEmail(email string) (*t.User, error)
	InsertUser(user t.User) (*t.User, error)
	FindAndUpdateSession(userID primitive.ObjectID, expiresAt time.Time) (*t.Session, error)
	CreateSession(session t.Session) (*t.Session, error)
}

type authRepository struct {
	userCollection    *mongo.Collection
	sessionCollection *mongo.Collection
}

func NewAuthRepository() AuthRepository {
	return &authRepository{
		userCollection:    db.Client.Database(db.DB_NAME).Collection("user"),
		sessionCollection: db.Client.Database(db.DB_NAME).Collection("session"),
	}
}

func (r *authRepository) FindUserByEmail(email string) (*t.User, error) {
	var user t.User
	err := r.userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) InsertUser(user t.User) (*t.User, error) {
	result, err := r.userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (r *authRepository) FindAndUpdateSession(userID primitive.ObjectID, expiresAt time.Time) (*t.Session, error) {
	var session t.Session
	update := bson.M{
		"$set": bson.M{"expires_at": expiresAt},
	}

	err := r.sessionCollection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"user_id": userID},
		update,
	).Decode(&session)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *authRepository) CreateSession(session t.Session) (*t.Session, error) {
	_, err := r.sessionCollection.InsertOne(context.TODO(), session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
