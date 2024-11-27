package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joshibbotson/gym-tracker-backend/internal/service"
)

type AuthHandler struct {
	Service service.AuthService
}

func (h *AuthHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		user, err := h.createUser(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")        // Allow only your frontend origin
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")          // Specify allowed methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Specify allowed headers
	w.Header().Set("Access-Control-Allow-Credentials", "true")                    // Allow credentials (cookies)

	// Handle preflight requests (OPTIONS method)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sessionInfo, err := h.login(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionInfo.SessionID,
		Expires:  sessionInfo.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Use true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	fmt.Printf("Response Headers: %v\n", w.Header())

	json.NewEncoder(w).Encode(sessionInfo)
}

func (h *AuthHandler) createUser(w http.ResponseWriter, r *http.Request) (*service.User, error) {
	body, err := getBody(r.Body)
	if err != nil {
		return nil, err
	}

	var user service.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	createdUser, err := h.Service.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) (*service.Session, error) {
	body, err := getBody(r.Body)
	if err != nil {
		return nil, err
	}
	type login struct {
		Email    string
		Password string
	}

	var loginDetails login
	if err := json.Unmarshal(body, &loginDetails); err != nil {
		return nil, err
	}

	println("email:", loginDetails.Email)
	println("password:", loginDetails.Password)

	session, err := h.Service.Login(loginDetails.Email, loginDetails.Password)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func getBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// func sessionMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Get the session cookie
// 		cookie, err := r.Cookie("session_id")
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		sessionID := cookie.Value

// 		// Check session in MongoDB
// 		var session service.Session
// 		err = sessionCollection.FindOne(context.TODO(), bson.M{"session_id": sessionID}).Decode(&session)
// 		if err != nil || session.ExpiresAt.Before(time.Now()) {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Attach user ID to the context for later use in the request lifecycle
// 		ctx := context.WithValue(r.Context(), "userID", session.UserID)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
