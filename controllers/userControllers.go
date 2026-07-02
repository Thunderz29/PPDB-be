package controllers

import (
	"ppdb-be/config"
	"ppdb-be/models"
	"ppdb-be/models/request"
	"ppdb-be/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Authenticate user
// @Description Authenticates a user with username and password, returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "Login Request Payload"
// @Success 200 {object} utils.Response{data=map[string]interface{}} "Success (contains token and user detail)"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /login [post]
func Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, utils.MsgBadRequest, err)
		return
	}

	var user models.User
	if err := config.DB.Where("user_name = ? AND is_deleted = ?", req.UserName, false).First(&user).Error; err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Username atau password salah", nil)
		return
	}

	if err := utils.ComparePassword(user.UserPassword, req.Password); err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Username atau password salah", nil)
		return
	}

	token, err := utils.GenerateToken(user.UserID, user.UserName, user.UserEmail)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
		"user":  user,
	})
}

// @Summary Register a new user
// @Description Register a new user in the system (Public endpoint).
// @Tags Users
// @Accept json
// @Produce json
// @Param body body request.CreateUserRequest true "Create User Payload"
// @Success 201 {object} utils.Response{data=models.User} "Created"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 409 {object} utils.ErrorResponse "Conflict"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /register [post]
func CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, utils.MsgBadRequest, err)
		return
	}

	var count int64
	config.DB.Model(&models.User{}).Where("(user_name = ? OR user_email = ?) AND is_deleted = ?", req.UserName, req.UserEmail, false).Count(&count)
	if count > 0 {
		utils.SendError(c, http.StatusConflict, "Username atau Email sudah terdaftar", nil)
		return
	}

	hashedPassword, err := utils.HashPassword(req.UserPassword)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	var creatorID *int64
	if id, exists := c.Get("userID"); exists {
		val := id.(int64)
		creatorID = &val
	}

	now := time.Now()
	user := models.User{
		RoleID:       req.RoleID,
		UserFullName: req.UserFullName,
		UserName:     req.UserName,
		UserPhone:    req.UserPhone,
		UserEmail:    req.UserEmail,
		UserPassword: hashedPassword,
		CreatedBy:    creatorID,
		CreatedOn:    &now,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "User berhasil dibuat", user)
}

// @Summary List all users
// @Description Retrieves a list of active users. Requires authentication token.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]models.User} "Success"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Where("is_deleted = ?", false).Find(&users).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Berhasil mengambil daftar user", users)
}

// @Summary Get user by ID
// @Description Retrieves details of a user by their user ID. Requires authentication token.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response{data=models.User} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	var user models.User
	if err := config.DB.Where("user_id = ? AND is_deleted = ?", id, false).First(&user).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, utils.MsgNotFound, nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Berhasil mengambil detail user", user)
}

// @Summary Update user details
// @Description Updates user details. Requires authentication token.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param body body request.UpdateUserRequest true "Update User Payload"
// @Success 200 {object} utils.Response{data=models.User} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, utils.MsgBadRequest, err)
		return
	}

	var user models.User
	if err := config.DB.Where("user_id = ? AND is_deleted = ?", id, false).First(&user).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, utils.MsgNotFound, nil)
		return
	}

	user.RoleID = req.RoleID
	user.UserFullName = req.UserFullName
	user.UserName = req.UserName
	user.UserPhone = req.UserPhone
	user.UserEmail = req.UserEmail

	if req.UserPassword != "" {
		hashed, err := utils.HashPassword(req.UserPassword)
		if err != nil {
			utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
			return
		}
		user.UserPassword = hashed
	}

	if mid, exists := c.Get("userID"); exists {
		val := mid.(int64)
		user.LastModifiedBy = &val
	}
	now := time.Now()
	user.LastModifiedOn = &now

	if err := config.DB.Save(&user).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "User berhasil diperbarui", user)
}

// @Summary Soft delete user
// @Description Soft deletes a user by setting is_deleted to true. Requires authentication token.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	var user models.User
	if err := config.DB.Where("user_id = ? AND is_deleted = ?", id, false).First(&user).Error; err != nil {
		utils.SendError(c, http.StatusNotFound, utils.MsgNotFound, nil)
		return
	}

	user.IsDeleted = true
	if mid, exists := c.Get("userID"); exists {
		val := mid.(int64)
		user.LastModifiedBy = &val
	}
	now := time.Now()
	user.LastModifiedOn = &now

	if err := config.DB.Save(&user).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "User berhasil dihapus", nil)
}
