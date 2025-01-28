package workout

// handle CRUD ops on workout configs
import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	m "github.com/joshibbotson/gym-tracker-backend/internal/middleware"
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
		m.PermissionMiddleware(h.handleCreateWorkout)(w, r)
	case http.MethodGet:
		path := r.URL.Path

		if path == "/workout/count" {
			m.PermissionMiddleware(h.handleReadActivitiesCount)(w, r)
			break
		}
		if dateParam := r.PathValue("date"); len(dateParam) > 0 {
			m.PermissionMiddleware(h.handleReadByDate)(w, r)
			break
		}
		m.PermissionMiddleware(h.handleReadActivites)(w, r)
	case http.MethodPatch:
		m.PermissionMiddleware(h.handleUpdateWorkout)(w, r)
	case http.MethodDelete:
		m.PermissionMiddleware(h.handleDeleteWorkout)(w, r)
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

	workout, err := h.Service.CreateWorkout(r.Context().Value("userID").(primitive.ObjectID), unmarshalledBody)
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

func (h *WorkoutHandler) handleReadByDate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(primitive.ObjectID)
	dateParam := r.PathValue("date")
	date, err := time.Parse("2006-01-02T15:04:05Z07:00", dateParam)
	if err != nil {
		http.Error(w, "Error parsing date", http.StatusInternalServerError)
	}
	workout, err := h.Service.GetWorkoutsByDate(userID, date)
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
	userID := r.Context().Value("userID").(primitive.ObjectID)

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

func (h *WorkoutHandler) handleReadActivitiesCount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(primitive.ObjectID)

	count, err := h.Service.GetActivityCountByUserId(userID)
	if err != nil {
		http.Error(w, "Error fetching activities count", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(count); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *WorkoutHandler) handleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	println("update handler hit")
	userID := r.Context().Value("userID").(primitive.ObjectID)

	body, getBodyErr := util.GetBody(r.Body)
	if getBodyErr != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	var unmarshalledBody t.UpdateWorkoutRequest
	err := json.Unmarshal(body, &unmarshalledBody)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println(unmarshalledBody)
	workouts, err := h.Service.UpdateWorkout(userID, unmarshalledBody)
	if err != nil {
		println(err)
		http.Error(w, "Error updating workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(workouts); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *WorkoutHandler) handleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	println("delete handler hit:", r.PathValue("id"))
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
