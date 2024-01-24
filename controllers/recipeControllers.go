package controllers

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"book-recipe-be-go/models/response"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func GetAllRecipes(c *gin.Context) {
	var recipes []models.Recipe
	var total int64

	minioClient, err := config.ConfigMinio()
    if err != nil {
        log.Fatalln(err)
        return
    }

	// filter
	pageNumber := c.DefaultQuery("pageNumber", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	recipeName := c.Query("recipeName")
	levelID := c.Query("levelId")
	categoryID := c.Query("categoryId")
	time := c.Query("time")
	sortBy := c.DefaultQuery("sortBy", "")

	// Mengonversi nilai string ke integer
	pageNumberInt, _ := strconv.Atoi(pageNumber)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	// Menetapkan nilai default jika tidak ada nilai atau nilai tidak valid
	if pageNumberInt <= 0 {
		pageNumberInt = 1
	}

	if pageSizeInt <= 0 {
		pageSizeInt = 10 // Atur ke ukuran halaman default
	}

	// Filter by recipe name
	db := config.DB.Model(&models.Recipe{})
	if recipeName != "" {
		db = db.Where("recipe_name LIKE ?", "%"+recipeName+"%")
	}

	// Filter by level id
	if levelID != "" {
		db = db.Where("level_id = ?", levelID)
	}

	// Filter by category id
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	// Filter by time cook
	if time != "" {
		timeCookInt, err := strconv.Atoi(time)

		if err != nil {
			response := response.MessageResponse{
				Message:    "Parameter time harus berupa angka",
				StatusCode: http.StatusBadRequest,
				Status:     "ERROR",
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if timeCookInt <= 30 {
			db = db.Where("time_cook <= ?", timeCookInt)
		} else if timeCookInt > 30 && timeCookInt <= 60 {
			db = db.Where("time_cook > 30 AND time_cook <= ?", timeCookInt)
		}

	}

	// Menentukan pengurutan berdasarkan sortBy
	if sortBy != "" {
		sortColumns := map[string]string{
			"recipeName": "LOWER(recipe_name)",
			"timeCook":   "time_cook",
		}

		sortParams := strings.Split(sortBy, ",")
		if len(sortParams) == 2 {
			columnName := sortParams[0]
			sortDirection := sortParams[1]

			if dbColumn, exists := sortColumns[columnName]; exists {
				switch sortDirection {
				case "asc":
					db = db.Order(fmt.Sprintf("%s ASC", dbColumn))
				case "desc":
					db = db.Order(fmt.Sprintf("%s DESC", dbColumn))
				default:
					// Jika arah pengurutan tidak valid, gunakan pengurutan default (ASC)
					db = db.Order(fmt.Sprintf("%s ASC", dbColumn))
				}
			}
		}
	}

	// Menggunakan Joins untuk mengambil data dari tabel levels dan categories
	if err := db.
		Preload("Category").
		Preload("Level").
		Order("LOWER(recipe_name) ASC").
		Where("is_deleted = ?", false).
		Count(&total).
		Limit(pageSizeInt).Offset((pageNumberInt - 1) * pageSizeInt).
		Find(&recipes).Error; err != nil {
		response := response.MessageResponse{
			Message:    "Terjadi kesalahan saat mengambil data resep",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Modify the response mapping
	var recipeEntries []response.RecipeEntry
	for _, recipe := range recipes {
		categoryName, err := getCategoryName(recipe.Category.CategoryID)
		if err != nil {
			log.Println("Error fetching category name:", err)
			categoryName = "Unknown Category"
		}

		levelName, err := getLevelName(recipe.Level.LevelID)
		if err != nil {
			log.Println("Error fetching level name:", err)
			levelName = "Unknown Level"
		}

		imageUrl, err := getImageURL(minioClient, recipe.ImageFilename)
		if err != nil {
			log.Println("Error getting image URL:", err)
			// Handle or return an appropriate response
			continue
		}

		entry := response.RecipeEntry{
			RecipeId:   recipe.RecipeID,
			Categories: response.CategoryInfo{CategoryId: recipe.CategoryID, CategoryName: categoryName},
			Levels:     response.LevelInfo{LevelId: recipe.LevelID, LevelName: levelName},
			RecipeName: recipe.RecipeName,
			ImageUrl:   imageUrl,
			Time:       recipe.TimeCook,
			IsFavorite: false, // You need to set this based on your logic
		}

		recipeEntries = append(recipeEntries, entry)
	}

	if len(recipeEntries) == 0 {
		response := response.MessageResponse{
			Message:    "Resep masakan tidak tersedia",
			StatusCode: http.StatusNotFound,
			Status:     "OK",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}
	
	response := response.RecipeListResponse{
		Total:      total,
		Data:       recipeEntries,
		Message:    "Berhasil memuat Resep Masakan Saya!",
		StatusCode: http.StatusOK,
		Status:     "Success",
	}
	c.JSON(http.StatusOK, response)
}

func getCategoryName(categoryID int) (string, error) {
	var category models.Category
	if err := config.DB.Model(&models.Category{}).Where("category_id = ?", categoryID).First(&category).Error; err != nil {
		return "", err
	}
	return category.CategoryName, nil
}

func getLevelName(levelID int) (string, error) {
	var level models.Level
	if err := config.DB.Model(&models.Level{}).Where("level_id = ?", levelID).First(&level).Error; err != nil {
		return "", err
	}
	return level.LevelName, nil
}

func getImageURL(minioClient *minio.Client, filename string) (string, error) {
	url, err := config.GetPublicLink(minioClient, filename)
	if err != nil {
		return "", err
	}

	return url, nil
}
