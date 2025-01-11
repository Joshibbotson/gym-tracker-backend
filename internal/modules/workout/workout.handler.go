package workout

// handle CRUD ops on workout configs
import (
	"encoding/json"
	"fmt"
	"net/http"

	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/workout/types"
	util "github.com/joshibbotson/gym-tracker-backend/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	case http.MethodGet:
		h.handleReadActivites(w, r)
	case http.MethodPatch:
		h.handleUpdateWorkout(w, r)
	case http.MethodDelete:
		h.handleDeleteWorkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WorkoutHandler) handleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	workout, err := h.Service.CreateWorkout(userID, unmarshalledBody)
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

func (h *WorkoutHandler) handleReadActivites(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workout, err := h.Service.GetWorkoutsByUserId(userID)
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

func (h *WorkoutHandler) handleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	println("update handler hit")
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, getBodyErr := util.GetBody(r.Body)
	if getBodyErr != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var unmarshalledBody t.UpdateWorkoutRequest
	err := json.Unmarshal(body, &unmarshalledBody)
	if err != nil {
		fmt.Println("err:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	workout, err := h.Service.UpdateWorkout(userID, unmarshalledBody)
	if err != nil {
		http.Error(w, "Error updating workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(workout); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *WorkoutHandler) handleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	println("delete handler hit")
	_, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idParam := r.PathValue("id")
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Convert the ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	success, err := h.Service.DeleteWorkout(objectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !success {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Workout deleted successfully"}`))
}
