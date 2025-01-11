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
		var allowedOrigin string
		if env == "production" {
			// In production, replace this with your frontend's HTTPS URL
			allowedOrigin = "https://gym-commits.netlify.app" // Replace with actual HTTPS frontend URL
		} else {
			// For local development, allow the local IP address with HTTP or HTTPS
			allowedOrigin = "http://localhost:4200, http://5.133.46.201:4200"
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests (OPTIONS method)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}
