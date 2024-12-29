package workout

// handle CRUD ops on workout configs
import (
	"net/http"
)

type WorkoutHandler struct {
	Service WorkoutService
}

func (h *WorkoutHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
	}
}
