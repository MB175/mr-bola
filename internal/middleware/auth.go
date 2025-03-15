package middleware

import (
	"errors"
	"github.com/MB175/mr-bola/pkg/utils"
	"net/http"
	"strings"
)

// Claims represents the JWT claims
type Claims struct {
	Subject string `json:"sub"`
}

// ExtractJWT extracts and validates the JWT token from request
func ExtractJWT(r *http.Request) (*Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("missing or invalid authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		return nil, err
	}

	return &Claims{Subject: claims.Subject}, nil
}
