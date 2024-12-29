package workout

// handle CRUD ops on workout configs
import (
	"encoding/json"
	"net/http"
)

type WorkoutHandler struct {
	Service WorkoutService
}

func NewWorkoutHandler(service WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		Service: service,
	}
}

// Handler for HTTP requests
func (h *WorkoutHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Handle POST request to create a workout
		h.handleCreateWorkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handle POST request to create a new workout
func (h *WorkoutHandler) handleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into the CreateWorkout struct
	var createWorkout CreateWorkout
	err := json.NewDecoder(r.Body).Decode(&createWorkout)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service to create a new workout
	workout, err := h.Service.createWorkout(createWorkout)
	if err != nil {
		http.Error(w, "Error creating workout", http.StatusInternalServerError)
		return
	}

	// Send a response with the created workout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(workout); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
