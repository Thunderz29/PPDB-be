package controllers

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"book-recipe-be-go/models/response"
	"book-recipe-be-go/utils"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetAllLevels mengembalikan semua kategori dengan hanya levelId dan levelName
func GetAllLevels(c *gin.Context) {
	var levels []models.Level
	result := config.DB.Select("level_id, level_name").Find(&levels)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, response.LevelListResponse{
			Message:    utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		})
		log.Println(utils.ErrInternalServer, result.Error.Error())
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
		Message:    utils.Success,
		StatusCode: http.StatusOK,
		Status:     "Success",
	})
	log.Printf(utils.Success)
}