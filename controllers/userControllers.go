package controllers

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"book-recipe-be-go/models/request"
	"book-recipe-be-go/models/response"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	if userExists := config.DB.Where("username = ?", request.Username).First(&models.User{}).Error; userExists == nil {
		response := response.MessageResponse{
			Message:    "Username telah digunakan oleh user yang telah mendaftar sebelumnya",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if request.Password != request.RetypePassword {
		response := response.MessageResponse{
			Message:    "Konfirmasi kata sandi tidak sama dengan kata sandi",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(request.Password) < 6 {
		response := response.MessageResponse{
			Message:    "Kata sandi tidak boleh kurang dari 6 karakter",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

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
		Status:     "OK",
	}

	c.JSON(http.StatusOK, response)
}

func SignIn(c *gin.Context) {
	var req request.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := response.MessageResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response := response.MessageResponse{
			Message:    "User not found",
			StatusCode: http.StatusNotFound,
			Status:     "ERROR",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Verifikasi password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		response := response.MessageResponse{
			Message:    "Invalid password",
			StatusCode: http.StatusUnauthorized,
			Status:     "ERROR",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		response := response.MessageResponse{
			Message:    "Failed to create token",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseData := response.SignInResponse{
		Data: response.UserData{
			ID:       user.UserID,
			Token:    tokenString,
			Type:     "Bearer",
			Username: user.Username,
			Role:     user.Role,
		},
		Message:    "Auth User Success",
		StatusCode: http.StatusOK,
		Status:     "OK",
	}

	c.JSON(http.StatusOK, responseData)
}
