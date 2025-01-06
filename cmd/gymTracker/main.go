package main

import (
	"context"
	"net/http"
	"regexp"
	"strings"

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
	http.HandleFunc("/workout/([^/]+)/parts/([0-9]+)", middlewareChain(workoutHandler.Handler))

	// put in env variable.
	http.ListenAndServe(":8888", nil)

}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

var routes = []route{}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method now allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}
