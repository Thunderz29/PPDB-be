package controllers

import (
	"net/http"
	"ppdb-be/config"
	"ppdb-be/models"
	"ppdb-be/models/request"
	"ppdb-be/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Where("is_deleted = ?", false).Find(&users).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, utils.MsgInternalServerError, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Berhasil mengambil daftar user", users)
}

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
