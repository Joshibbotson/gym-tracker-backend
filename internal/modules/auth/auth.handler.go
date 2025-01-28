package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joshibbotson/gym-tracker-backend/internal/modules/auth/constants"
	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/auth/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	Service AuthService
}

var (
	googleOauthConfig = oauth2.Config{
		ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8888/auth/google/callback"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	oauthStateString = generateStateString(32)
)

// getEnv fetches the value of an environment variable or returns a default value if not set
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func generateStateString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Failed to generate random state string: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

/*
* Need  a better setup for auth, this is messy, login and user creation should all just belong to auth
to simplify things.
*/
// func (h *AuthHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodPost:
// 		user, err := h.createUser(w, r)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(user)

// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func (h *AuthHandler) HandleUserDetails(w http.ResponseWriter, r *http.Request) {
// 	userId, err :=
// }

// handleGoogleLogin redirects the user to Google's OAuth 2.0 server
func (h *AuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL("random_state_string", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *AuthHandler) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "State not found", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := googleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the token to get user information
	client := googleOauthConfig.Client(r.Context(), token)
	userInfoResponse, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer userInfoResponse.Body.Close()

	// Parse and display the user information
	var userInfo t.AuthData
	if err := json.NewDecoder(userInfoResponse.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	userInfo.AuthProvider = constants.AuthProvidersGoogle

	sessionInfo, err := h.Service.LoginOrCreateUser(userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	env := os.Getenv("GO_ENV")
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionInfo.SessionID,
		Expires:  sessionInfo.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   (env == "production"),
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	REDIRECT_URL := os.Getenv("REDIRECT_URL")

	// Redirect before any body content is written
	redirectURL := fmt.Sprintf(REDIRECT_URL+"/redirect-auth/?name=%s&email=%s", sessionInfo.Name, sessionInfo.Email)
	// Perform the redirect to the frontend
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("GO_ENV")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   (env == "production"),
		SameSite: http.SameSiteLaxMode})

	w.WriteHeader(http.StatusOK)
}

// func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	// authenticate the user
// 	sessionInfo, err := h.login(w, r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	env := os.Getenv("GO_ENV")
// 	cookie := &http.Cookie{
// 		Name:     "session_token",
// 		Value:    sessionInfo.SessionID,
// 		Expires:  sessionInfo.ExpiresAt,
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   (env == "production"),
// 		SameSite: http.SameSiteLaxMode,
// 	}
// 	http.SetCookie(w, cookie)

// 	w.WriteHeader(http.StatusOK)

// 	type UserDetails struct {
// 		Name  string `json:"name"`
// 		Email string `json:"email"`
// 	}
// 	userDetails := UserDetails{
// 		Name:  sessionInfo.Name,
// 		Email: sessionInfo.Email,
// 	}

// 	json.NewEncoder(w).Encode(userDetails)

// }

// func (h *AuthHandler) createUser(w http.ResponseWriter, r *http.Request) (*t.User, error) {
// 	body, err := u.GetBody(r.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var user t.User
// 	if err := json.Unmarshal(body, &user); err != nil {
// 		return nil, err
// 	}

// 	createdUser, err := h.Service.CreateLocalUser(user.Name, user.Email, user.Password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return createdUser, nil
// }

// func (h *AuthHandler) login(_ http.ResponseWriter, r *http.Request) (*t.Session, error) {
// 	body, err := u.GetBody(r.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	type login struct {
// 		Email    string
// 		Password string
// 	}

// 	var loginDetails login
// 	if err := json.Unmarshal(body, &loginDetails); err != nil {
// 		return nil, err
// 	}

// 	println("email:", loginDetails.Email)
// 	println("password:", loginDetails.Password)

// 	session, err := h.Service.Login(loginDetails.Email, loginDetails.Password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return session, nil
// }
