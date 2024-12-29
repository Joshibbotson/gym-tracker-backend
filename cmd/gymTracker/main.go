package main

import (
	"net/http"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
	"github.com/joshibbotson/gym-tracker-backend/internal/middleware"
	"github.com/joshibbotson/gym-tracker-backend/internal/modules/auth"
	"github.com/joshibbotson/gym-tracker-backend/internal/modules/workout"
)

func main() {
	db.ConnectDB()
	defer db.DisconnectDB()

	authService := auth.NewAuthService()
	authHandler := &auth.AuthHandler{Service: authService}

	workoutService := workout.NewWorkoutService()
	workoutHandler := &workout.WorkoutHandler{Service: workoutService}

	http.HandleFunc("/auth", authHandler.Handler)
	http.HandleFunc("/auth/login", authHandler.LoginHandler)
	http.Handle("/workout", middleware.SessionMiddleware(http.HandlerFunc(workoutHandler.Handler)))

	// put in env variable.
	http.ListenAndServe(":8888", nil)

}
