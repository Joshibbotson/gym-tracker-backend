package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
	"github.com/joshibbotson/gym-tracker-backend/internal/modules/auth"
	"go.mongodb.org/mongo-driver/bson"
)

const DB_NAME = "gym-tracker"

func getUserBySessionId(sessionId string) (auth.Session, error) {
	sessionCollection := db.Client.Database(DB_NAME).Collection("session")

	var session auth.Session
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := sessionCollection.FindOne(ctx, bson.M{"session_id": sessionId}).Decode(&session)
	if err != nil {
		return auth.Session{}, fmt.Errorf("failed to fetch session: %v", err)
	}

	return session, nil
}

func SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(2)

		// Get the session cookie

		cookie, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("Cookie error:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value
		fmt.Printf("SessionID: %s\n", sessionID)

		// Validate session using the AuthService
		session, err := getUserBySessionId(sessionID)
		if err != nil {
			fmt.Printf("Session fetch error: %v\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			fmt.Println("Session expired")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Attach user ID to the context for later use in the request lifecycle
		fmt.Printf("Session valid for UserID: %s\n", session.UserID)
		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
