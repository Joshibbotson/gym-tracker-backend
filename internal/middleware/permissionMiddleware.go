package middleware

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PermissionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, ok := r.Context().Value("userID").(primitive.ObjectID)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
