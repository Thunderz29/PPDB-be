package utils

import (
    "book-recipe-be-go/config"
    "book-recipe-be-go/models"
)

func GetFullnameByUserID(userID uint) (string, error) {
	var user models.User
	if err := config.DB.Model(&models.User{}).Where("user_id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Fullname, nil
}