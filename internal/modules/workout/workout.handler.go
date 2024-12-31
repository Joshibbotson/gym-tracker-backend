package workout

// handle CRUD ops on workout configs
import (
	"encoding/json"
	"fmt"
	"net/http"

	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/workout/types"
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
	body, getBodyErr := util.GetBody(r.Body)
	if getBodyErr != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var unmarshalledBody t.CreateWorkoutRequest
	err := json.Unmarshal(body, &unmarshalledBody)
	if err != nil {
		fmt.Println("err:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	workout, err := h.Service.createWorkout(unmarshalledBody)
	if err != nil {
		http.Error(w, "Error creating workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(workout); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
