package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing JWT tokens (for demo purposes only)
var secretKey = []byte("supersecretkey")

// GenerateJWT creates a new JWT token for a given username
func GenerateJWT(username string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // 1-hour expiry
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT validates a given JWT token
func ValidateJWT(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
