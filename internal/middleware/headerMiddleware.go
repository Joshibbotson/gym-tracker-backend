package middleware

import (
	"fmt"
	"net/http"
	"os"
)

func HeaderMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req cookies:", r.Cookies())

		env := os.Getenv("GO_ENV")

		// Set allowed origins based on the environment
		var allowedOrigins []string
		if env == "production" {
			allowedOrigins = []string{"https://gym-tracker.joshibbotson.com"}
		} else if env == "staging" {
			allowedOrigins = []string{"https://gym-tracker-staging.joshibbotson.com"}
		} else {
			allowedOrigins = []string{"http://localhost:4200", "http://5.133.46.201:4200"}
		}

		origin := r.Header.Get("Origin")
		if contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests (OPTIONS method)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func contains(origins []string, target string) bool {
	for _, o := range origins {
		if o == target {
			return true
		}
	}
	return false
}
