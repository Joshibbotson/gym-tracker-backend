package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/joshibbotson/gym-tracker-backend/internal/service"
)

type AuthHandler struct {
	service service.AuthRepository
}

func (h *AuthHandler) handler(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) createUser(w http.ResponseWriter, r *http.Request) (*service.User, error) {
	body, err := getBody(r.Body)
	if err != nil {
		return nil, err
	}

	var user service.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	createdUser, err := h.service.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func getBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
