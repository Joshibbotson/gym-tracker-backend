package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/joshibbotson/gym-tracker-backend/internal/modules/auth"
)

// Middleware struct to encapsulate dependencies
type Middleware struct {
	AuthService *auth.AuthService
}

// NewMiddleware constructor
func NewMiddleware(authService *auth.AuthService) *Middleware {
	return &Middleware{AuthService: authService}
}

// SessionMiddleware validates the session
func (m *Middleware) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value

		// Validate session using the AuthService
		session, err := m.AuthService.GetUserBySessionId(sessionID)
		if err != nil || session.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Attach user ID to the context for later use in the request lifecycle
		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
