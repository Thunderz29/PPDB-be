package controllers

import (
	"net/http"
	"ppdb-be/config"
	"ppdb-be/models"
	"ppdb-be/models/request"
	"ppdb-be/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func parseUserAgent(ua string) (platform, browser, device string) {
	if ua == "" {
		return "Unknown", "Unknown", "Unknown"
	}
	lowerUA := strings.ToLower(ua)

	if strings.Contains(lowerUA, "edg") {
		browser = "Edge"
	} else if strings.Contains(lowerUA, "chrome") {
		browser = "Chrome"
	} else if strings.Contains(lowerUA, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(lowerUA, "safari") {
		browser = "Safari"
	} else {
		browser = "Unknown"
	}

	if strings.Contains(lowerUA, "windows") {
		platform = "Windows"
		device = "PC"
	} else if strings.Contains(lowerUA, "macintosh") || strings.Contains(lowerUA, "mac os") {
		platform = "macOS"
		device = "Mac"
	} else if strings.Contains(lowerUA, "android") {
		platform = "Android"
		device = "Mobile"
	} else if strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "ipad") {
		platform = "iOS"
		if strings.Contains(lowerUA, "ipad") {
			device = "iPad"
		} else {
			device = "iPhone"
		}
	} else if strings.Contains(lowerUA, "linux") {
		platform = "Linux"
		device = "PC"
	} else {
		platform = "Unknown"
		device = "Unknown"
	}

	return platform, browser, device
}

// @Summary Authenticate user and record session
// @Description Authenticates a user with username and password, returns access token (1 day), refresh token (7 days), and records the session.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "Login Request Payload"
// @Success 200 {object} utils.Response{data=map[string]interface{}} "Success"
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

	accessToken, err := utils.GenerateToken(user.UserID, user.UserName, user.UserEmail, 24*time.Hour)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	refreshToken, err := utils.GenerateToken(user.UserID, user.UserName, user.UserEmail, 7*24*time.Hour)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	ua := c.Request.UserAgent()
	platform, browser, device := parseUserAgent(ua)
	ip := c.ClientIP()
	now := time.Now()
	expiredAt := now.Add(7 * 24 * time.Hour)

	config.DB.WithContext(c).Model(&models.UserSession{}).
		Where("user_id = ? AND is_active = ?", user.UserID, true).
		Updates(map[string]interface{}{
			"is_active":  false,
			"is_revoked": true,
			"logout_at":  &now,
		})

	session := models.UserSession{
		UserID:       user.UserID,
		AccessToken:  accessToken,
		RefreshToken: &refreshToken,
		LoginAt:      now,
		ExpiredAt:    expiredAt,
		IPAddress:    &ip,
		UserAgent:    &ua,
		DeviceName:   &device,
		Platform:     &platform,
		Browser:      &browser,
		IsActive:     true,
		IsRevoked:    false,
		CreatedOn:    now,
	}

	if err := config.DB.WithContext(c).Create(&session).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Login berhasil", gin.H{
		"session_id":    session.SessionID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
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

	user := models.User{
		RoleID:       req.RoleID,
		UserFullName: req.UserFullName,
		UserName:     req.UserName,
		UserPhone:    req.UserPhone,
		UserEmail:    req.UserEmail,
		UserPassword: hashedPassword,
	}

	if err := config.DB.WithContext(c).Create(&user).Error; err != nil {
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

	if err := config.DB.WithContext(c).Save(&user).Error; err != nil {
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

	if err := config.DB.WithContext(c).Save(&user).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "User berhasil dihapus", nil)
}

// @Summary Logout user and invalidate session
// @Description Logs out the user by revoking the active session associated with the Bearer token.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Success"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /logout [post]
func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.SendError(c, http.StatusBadRequest, utils.MsgBadRequest, "Authorization header is required")
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.SendError(c, http.StatusBadRequest, utils.MsgBadRequest, "Invalid token format")
		return
	}
	tokenStr := parts[1]

	var session models.UserSession
	if err := config.DB.WithContext(c).Where("access_token = ? AND is_active = ?", tokenStr, true).First(&session).Error; err != nil {
		utils.SendSuccess(c, http.StatusOK, "Logout berhasil", nil)
		return
	}

	now := time.Now()
	session.IsActive = false
	session.IsRevoked = true
	session.LogoutAt = &now

	if err := config.DB.WithContext(c).Save(&session).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Logout berhasil", nil)
}

