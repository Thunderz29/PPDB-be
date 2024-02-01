package controllers

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"book-recipe-be-go/models/request"
	"book-recipe-be-go/models/response"
	"book-recipe-be-go/utils"
	"fmt"
	"log"
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
		log.Println(utils.ErrBadRequest, err.Error())
		return
	}

	if userExists := config.DB.Where("username = ?", request.Username).First(&models.User{}).Error; userExists == nil {
		response := response.MessageResponse{
			Message:    utils.ErrUserAlreadyExist,
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		log.Println(utils.ErrBadRequest, utils.ErrUserAlreadyExist)
		return
	}

	if request.Password != request.RetypePassword {
		response := response.MessageResponse{
			Message:    utils.ErrConfirmPassword,
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		log.Println(utils.ErrBadRequest, utils.ErrConfirmPassword)
		return
	}

	if len(request.Password) < 6 {
		response := response.MessageResponse{
			Message:    utils.ErrValidatePassword,
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		log.Println(utils.ErrBadRequest, utils.ErrValidatePassword)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		response := response.MessageResponse{
			Message:    utils.ErrHashPassword,
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		log.Println(utils.ErrInternalServer, utils.ErrHashPassword)
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
			Message:    utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		log.Println(utils.ErrInternalServer, err.Error())
		return
	}

	response := response.MessageResponse{
		Message:   fmt.Sprintf(utils.SuccSignUp, request.Username),
		StatusCode: http.StatusOK,
		Status:     "OK",
	}
	c.JSON(http.StatusOK, response)
	log.Printf(response.Message)
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
		log.Println(utils.ErrBadRequest, err.Error())
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response := response.MessageResponse{
			Message:    utils.ErrUserNotFound,
			StatusCode: http.StatusNotFound,
			Status:     "ERROR",
		}
		c.JSON(http.StatusUnauthorized, response)
		log.Println(utils.ErrUserNotFound, err.Error())
		return
	}

	// Verifikasi password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		response := response.MessageResponse{
			Message:    utils.ErrInvalidPassword,
			StatusCode: http.StatusUnauthorized,
			Status:     "ERROR",
		}
		c.JSON(http.StatusUnauthorized, response)
		log.Println(utils.ErrInvalidPassword, err.Error())
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
			Message:    utils.ErrCreateToken,
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		log.Println(utils.ErrInternalServer, utils.ErrCreateToken)
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
		Message:    utils.SuccSignIn,
		StatusCode: http.StatusOK,
		Status:     "OK",
	}

	c.JSON(http.StatusOK, responseData)
	log.Printf(responseData.Message)
}
