package auth

import (
	"encoding/json"
	"io"
	"net/http"
)

type AuthHandler struct {
	Service AuthService
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
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight requests (OPTIONS method)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// authenticate the user
	sessionInfo, err := h.login(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionInfo.SessionID,
		Expires:  sessionInfo.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	})

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

func (h *AuthHandler) createUser(w http.ResponseWriter, r *http.Request) (*User, error) {
	body, err := getBody(r.Body)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	createdUser, err := h.Service.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) (*Session, error) {
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
