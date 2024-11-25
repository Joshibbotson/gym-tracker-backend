package main

import (
	"net/http"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
	"github.com/joshibbotson/gym-tracker-backend/internal/handler"
	"github.com/joshibbotson/gym-tracker-backend/internal/service"
)

func main() {
	db.ConnectDB()
	defer db.DisconnectDB()

	authService := service.NewAuthService()
	authHandler := &handler.AuthHandler{Service: authService}
	http.HandleFunc("/auth", authHandler.Handler)
	http.HandleFunc("/auth/login", authHandler.LoginHandler)

	// put in env variable.
	http.ListenAndServe(":8888", nil)

}
