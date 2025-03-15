package auth

import (
	"encoding/json"
	"github.com/MB175/mr-bola/pkg/utils"
	"net/http"
)

// AuthRequest represents the authentication request payload
type AuthRequest struct {
	Username string `json:"username"`
}

// AuthResponse represents the JWT response
type AuthResponse struct {
	Token string `json:"token"`
}

// AuthHandler handles mock authentication and returns a JWT token
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Username == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := utils.GenerateJWT(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}
