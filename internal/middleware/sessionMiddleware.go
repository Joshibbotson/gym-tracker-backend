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
	err := sessionCollection.FindOne(context.TODO(), bson.M{"session_id": sessionId}).Decode(&session)
	if err != nil {
		return auth.Session{}, fmt.Errorf("failed to fetch session: %v", err)
	}

	return session, nil
}

// SessionMiddleware validates the session
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value

		// Validate session using the AuthService
		session, err := getUserBySessionId(sessionID)
		if err != nil || session.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Attach user ID to the context for later use in the request lifecycle
		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
