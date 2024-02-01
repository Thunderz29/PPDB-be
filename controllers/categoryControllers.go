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

// GetAllCategories mengembalikan semua kategori dengan hanya categoryId dan categoryName
func GetAllCategories(c *gin.Context) {
	var categories []models.Category
	result := config.DB.Select("category_id, category_name").Find(&categories)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, response.CategoryListResponse{
			Message:    utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		})
		log.Println(utils.ErrInternalServer, result.Error.Error())
		return
	}

	// Mengonversi model Category ke CategoryInfo dalam respons
	var categoryInfos []response.CategoryInfo
	for _, category := range categories {
		categoryInfos = append(categoryInfos, response.CategoryInfo{
			CategoryId:   category.CategoryID,
			CategoryName: category.CategoryName,
		})
	}

	c.JSON(http.StatusOK, response.CategoryListResponse{
		Data:       categoryInfos,
		Message:    utils.Success,
		StatusCode: http.StatusOK,
		Status:     "Success",
	})
	log.Printf(utils.Success)
}