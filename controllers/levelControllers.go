package controllers

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"book-recipe-be-go/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAllLevels mengembalikan semua kategori dengan hanya levelId dan levelName
func GetAllLevels(c *gin.Context) {
	var levels []models.Level
	result := config.DB.Select("level_id, level_name").Find(&levels)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, response.LevelListResponse{
			Message:    "Terjadi kesalahan server. Silakan coba kembali",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		})
		return
	}

	// Mengonversi model Level ke LevelInfo dalam respons
	var levelInfos []response.LevelInfo
	for _, level := range levels {
		levelInfos = append(levelInfos, response.LevelInfo{
			LevelId:   level.LevelID,
			LevelName: level.LevelName,
		})
	}

	c.JSON(http.StatusOK, response.LevelListResponse{
		Data:       levelInfos,
		Message:    "Pesan Sukses",
		StatusCode: http.StatusOK,
		Status:     "Success",
	})
}