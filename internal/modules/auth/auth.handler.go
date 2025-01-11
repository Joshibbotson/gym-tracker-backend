package auth

import (
	"encoding/json"
	"net/http"
	"os"

	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/auth/types"

	u "github.com/joshibbotson/gym-tracker-backend/internal/util"
)

type AuthHandler struct {
	Service AuthService
}

func (h *AuthHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
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
	// authenticate the user
	sessionInfo, err := h.login(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	env := os.Getenv("GO_ENV")

	if env == "production" {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionInfo.SessionID,
			Expires:  sessionInfo.ExpiresAt,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionInfo.SessionID,
			Expires:  sessionInfo.ExpiresAt,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
	}

	w.WriteHeader(http.StatusOK)

	type UserDetails struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	userDetails := UserDetails{
		Name:  sessionInfo.Name,
		Email: sessionInfo.Email,
	}

	json.NewEncoder(w).Encode(userDetails)

}

func (h *AuthHandler) createUser(w http.ResponseWriter, r *http.Request) (*t.User, error) {
	body, err := u.GetBody(r.Body)
	if err != nil {
		return nil, err
	}

	var user t.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	createdUser, err := h.Service.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (h *AuthHandler) login(_ http.ResponseWriter, r *http.Request) (*t.Session, error) {
	body, err := u.GetBody(r.Body)
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
