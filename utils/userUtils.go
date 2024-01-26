package utils

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func GetFullnameByUserID(userID uint) (string, error) {
	var user models.User
	if err := config.DB.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Fullname, nil
}

func GetUserIdFromToken(tokenString string) (int, error) {
	// Load SECRET from environment variable
	secretKey := os.Getenv("SECRET")
    log.Println("Secret Key:", secretKey)
    log.Println("Received Token:", tokenString)

	// Hapus "Bearer " dari awalan token
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
        log.Println("Error parsing token:", err)
		return 0, err
	}

	// Verifikasi token
	if !token.Valid {
        log.Println("Token is not valid")
		return 0, errors.New("Token tidak valid")
	}

	// Ambil nilai subject (sub) dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
        log.Println("Failed to get token claims")
		return 0, errors.New("Gagal mendapatkan klaim token")
	}

	subject, ok := claims["sub"].(float64)
	if !ok {
        log.Println("Claim 'sub' does not contain a valid subject")
		return 0, errors.New("Klaim token tidak berisi subjek yang valid")
	}

	return int(subject), nil
}



func GetusernameByUserID(userID uint) (string, error) {
	var user models.User
	if err := config.DB.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Username, nil
}