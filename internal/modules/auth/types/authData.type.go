package types

import (
	a "github.com/joshibbotson/gym-tracker-backend/internal/modules/auth/constants"
)

type AuthData struct {
	Email         string          `json:"email"`
	Surname       string          `json:"family_name"`
	FirstName     string          `json:"given_name"`
	AuthId        string          `json:"id"`
	Name          string          `json:"name"`
	PictureUrl    string          `json:"picture"`
	VerifiedEmail bool            `json:"verified_email"`
	AuthProvider  a.AuthProviders `json:"authProvider,omitempty" bson:"authProvider,omitempty"`
}
