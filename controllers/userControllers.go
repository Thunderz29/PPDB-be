package controllers

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"book-recipe-be-go/models/request"
	"book-recipe-be-go/models/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var request request.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response := response.MessageResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Cek apakah username sudah ada
	if userExists := config.DB.Where("username = ?", request.Username).First(&models.User{}).Error; userExists == nil {
		response := response.MessageResponse{
			Message:    "Username telah digunakan oleh user yang telah mendaftar sebelumnya",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Periksa apakah password cocok dengan retype password
	if request.Password != request.RetypePassword {
		response := response.MessageResponse{
			Message:    "Konfirmasi kata sandi tidak sama dengan kata sandi",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validasi panjang password
	if len(request.Password) < 6 {
		response := response.MessageResponse{
			Message:    "Kata sandi tidak boleh kurang dari 6 karakter",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		response := response.MessageResponse{
			Message:    "Error Hashing Password!",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	user := models.User{
		Username:    request.Username,
		Fullname:    request.Fullname,
		Password:    string(hashedPassword),
		Role:        "User",
		IsDeleted:   false,
		CreatedBy:   request.Username, 
		CreatedTime: time.Now(),
		ModifiedTime: time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		response := response.MessageResponse{
			Message:    "Terjadi kesalahan server. Silakan coba kembali",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.MessageResponse{
		Message:   "User " + request.Username + " registered successfully!",
		StatusCode: http.StatusOK,
		Status:     "Success",
	}

	c.JSON(http.StatusOK, response)
}
