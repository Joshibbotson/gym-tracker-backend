package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
	m "github.com/joshibbotson/gym-tracker-backend/internal/middleware"
	"github.com/joshibbotson/gym-tracker-backend/internal/modules/auth"
	"github.com/joshibbotson/gym-tracker-backend/internal/modules/workout"
)

func main() {
	println("initting main")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}
	println("port:", port)

	db.ConnectDB()
	defer db.DisconnectDB()
	middlewareChain := m.MiddlewareChain(m.HeaderMiddleware, m.SessionMiddleware)

	authService := auth.NewAuthService()
	authHandler := &auth.AuthHandler{Service: authService}

	workoutService := workout.NewWorkoutService()
	workoutHandler := &workout.WorkoutHandler{Service: workoutService}

	http.HandleFunc("/auth", authHandler.UserHandler)
	http.HandleFunc("/auth/login", m.HeaderMiddleware(authHandler.LoginHandler))
	http.HandleFunc("/auth/google/login", m.HeaderMiddleware((authHandler.HandleGoogleLogin)))
	http.HandleFunc("/auth/google/callback", m.HeaderMiddleware((authHandler.HandleOAuth2Callback)))
	http.HandleFunc("/workout", middlewareChain(workoutHandler.Handler))
	http.HandleFunc("/workout/{id}", middlewareChain(workoutHandler.Handler))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	// put in env variable.
	http.ListenAndServe("0.0.0.0:"+port, nil)

}
