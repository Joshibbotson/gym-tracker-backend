package workout

// handle CRUD ops on workout configs
import (
	"encoding/json"
	"fmt"
	"net/http"

	util "github.com/joshibbotson/gym-tracker-backend/internal/util"
)

type WorkoutHandler struct {
	Service WorkoutService
}

func NewWorkoutHandler(service WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		Service: service,
	}
}

func (h *WorkoutHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreateWorkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WorkoutHandler) handleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into the CreateWorkout struct

	body, getBodyErr := util.GetBody(r.Body)
	if getBodyErr != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var decodedBody CreateWorkoutRequest
	err := json.Unmarshal(body, &decodedBody)
	if err != nil {
		fmt.Println("err:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service to create a new workout
	workout, err := h.Service.createWorkout(decodedBody)
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
