package controllers

import (
	"net/http"
	"ppdb-be/models/request"
	"ppdb-be/services"
	"ppdb-be/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Login(c *gin.Context)
	CreateUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	Logout(c *gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

// @Summary Authenticate user and record session
// @Description Authenticates a user with email and password, returns access token (1 day), refresh token (7 days), and records the session.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "Login Request Payload"
// @Success 200 {object} utils.Response{data=map[string]interface{}} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /login [post]
func (ctrl *userController) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, utils.MsgBadRequest, err)
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	result, err := ctrl.userService.Login(c.Request.Context(), req.Email, req.Password, ip, ua)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Login berhasil", result)
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
func (ctrl *userController) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, utils.MsgBadRequest, err)
		return
	}

	user, err := ctrl.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "sudah terdaftar") {
			utils.SendError(c, http.StatusConflict, err.Error(), nil)
		} else {
			utils.SendError(c, http.StatusInternalServerError, err.Error(), nil)
		}
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
func (ctrl *userController) GetUsers(c *gin.Context) {
	users, err := ctrl.userService.GetUsers(c.Request.Context())
	if err != nil {
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
func (ctrl *userController) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	user, err := ctrl.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
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
func (ctrl *userController) UpdateUser(c *gin.Context) {
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

	user, err := ctrl.userService.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "User tidak ditemukan" {
			utils.SendError(c, http.StatusNotFound, utils.MsgNotFound, nil)
		} else {
			utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		}
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
func (ctrl *userController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	if err := ctrl.userService.DeleteUser(c.Request.Context(), id); err != nil {
		if err.Error() == "User tidak ditemukan" {
			utils.SendError(c, http.StatusNotFound, utils.MsgNotFound, nil)
		} else {
			utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		}
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
func (ctrl *userController) Logout(c *gin.Context) {
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

	if err := ctrl.userService.Logout(c.Request.Context(), tokenStr); err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Logout berhasil", nil)
}
