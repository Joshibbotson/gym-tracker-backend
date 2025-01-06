package main

import (
	"net/http"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
	m "github.com/joshibbotson/gym-tracker-backend/internal/middleware"
	"github.com/joshibbotson/gym-tracker-backend/internal/modules/auth"
	"github.com/joshibbotson/gym-tracker-backend/internal/modules/workout"
)

func main() {

	db.ConnectDB()
	defer db.DisconnectDB()
	middlewareChain := m.MiddlewareChain(m.HeaderMiddleware, m.SessionMiddleware)

	authService := auth.NewAuthService()
	authHandler := &auth.AuthHandler{Service: authService}

	workoutService := workout.NewWorkoutService()
	workoutHandler := &workout.WorkoutHandler{Service: workoutService}

	http.HandleFunc("/auth", authHandler.UserHandler)
	http.HandleFunc("/auth/login", m.HeaderMiddleware(authHandler.LoginHandler))
	http.HandleFunc("/workout", middlewareChain(workoutHandler.Handler))
	http.HandleFunc("/workout/{id}", middlewareChain(workoutHandler.Handler))

	// put in env variable.
	http.ListenAndServe(":8888", nil)

}
