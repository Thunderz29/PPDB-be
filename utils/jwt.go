package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("SECRET"))

// Claims defines custom JWT claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"user_email"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a specific user
func GenerateToken(userID int64, username, email string) (string, error) {
	// Fallback key if SECRET is not set in environment
	if len(jwtKey) == 0 {
		jwtKey = []byte("default_jwt_secret_key_change_me")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		UserName: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateToken validates a JWT token and returns claims
func ValidateToken(tokenString string) (*Claims, error) {
	if len(jwtKey) == 0 {
		jwtKey = []byte("default_jwt_secret_key_change_me")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
