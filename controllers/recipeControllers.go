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
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ToggleFavorite(c *gin.Context) {
	recipeID := c.Param("recipeId")
	var toggleRequest request.ToggleFavoriteRequest
	var favoriteFood models.FavoriteFood

	if err := c.ShouldBindJSON(&toggleRequest); err != nil {
		response := response.MessageResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Mengambil fullname berdasarkan userId
	fullname, err := utils.GetFullnameByUserID(toggleRequest.UserId)
	if err != nil {
		response := response.MessageResponse{
			Message:    "Error getting user information",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Mengkonversi recipId menjadi integer
	recipeIDInt, err := strconv.Atoi(recipeID)
	if err != nil {
		response := response.DataResponse{
			Total: 0,
			Data: nil,
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Mengambil recipe name berdasarkan recipeId
	recipeName, err := utils.GetRecipeNameByRecipeID(recipeIDInt)
	if err != nil {
		response := response.DataResponse{
			Total: 0,
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Check if recipe is already a favorite
	result := config.DB.Where("user_id = ? AND recipe_id = ?", toggleRequest.UserId, recipeIDInt).First(&favoriteFood)
	if result.Error == nil {
		// Recipe is already a favorite, toggle the is_favorite field
		favoriteFood.IsFavorite = !favoriteFood.IsFavorite
		favoriteFood.ModifiedBy = fullname
		favoriteFood.ModifiedTime = time.Now()

		if err := config.DB.Save(&favoriteFood).Error; err != nil {
			response := response.DataResponse{
				Total: 0,
				Data: nil,
				Message:    "Terjadi kesalahan server. Silakan coba kembali.",
				StatusCode: http.StatusInternalServerError,
				Status:     "ERROR",
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		message := "Resep " + recipeName + " berhasil ditambahkan ke dalam favorit"
		if !favoriteFood.IsFavorite {
			message = "Resep " + recipeName + " berhasil dihapus dari favorit"
		}

		response := response.DataResponse{
			Total: 1,
			Data: nil,
			Message:    message,
			StatusCode: http.StatusOK,
			Status:     "OK",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// Recipe is not yet a favorite, insert a new favorite record
	newFavorite := models.FavoriteFood{
		UserID:       int(toggleRequest.UserId),
		RecipeID:     recipeIDInt,
		IsFavorite:   true,
		CreatedBy:    fullname,
		CreatedTime:  time.Now(),
		ModifiedTime: time.Now(),
	}

	if err := config.DB.Create(&newFavorite).Error; err != nil {
		response := response.DataResponse{
			Total: 0,
			Data: nil,
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.DataResponse{
		Total: 1,
		Data: nil,
		Message:    "Resep " + recipeName + "berhasil ditambahkan ke dalam favorit",
		StatusCode: http.StatusOK,
		Status:     "OK",
	}
	c.JSON(http.StatusOK, response)
}

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
	userID := c.Query("userId")

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

	if(userID == ""){
		response := response.MessageResponse{
			Message:    "User tidak valid",
			StatusCode: http.StatusUnauthorized,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
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
				Message:    err.Error(),
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
		categoryName, err := utils.GetCategoryName(recipe.Category.CategoryID)
		if err != nil {
			log.Println("Error fetching category name:", err)
			categoryName = "Unknown Category"
		}

		levelName, err := utils.GetLevelName(recipe.Level.LevelID)
		if err != nil {
			log.Println("Error fetching level name:", err)
			levelName = "Unknown Level"
		}

		imageUrl, err := utils.GetImageURL(minioClient, recipe.ImageFilename)
		if err != nil {
			log.Println("Error getting image URL:", err)
			continue
		}

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			response := response.MessageResponse{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Status:     "ERROR",
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		isFavorite := utils.CheckFavoriteRecipe(uint(userIDInt), recipe.RecipeID)

		entry := response.RecipeEntry{
			RecipeId:   recipe.RecipeID,
			Categories: response.CategoryInfo{CategoryId: recipe.CategoryID, CategoryName: categoryName},
			Levels:     response.LevelInfo{LevelId: recipe.LevelID, LevelName: levelName},
			RecipeName: recipe.RecipeName,
			ImageUrl:   imageUrl,
			Time:       recipe.TimeCook,
			IsFavorite: isFavorite,
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
	
	response := response.DataResponse{
		Total:      total,
		Data:       recipeEntries,
		Message:    "Berhasil memuat Resep Masakan Saya!",
		StatusCode: http.StatusOK,
		Status:     "Success",
	}
	c.JSON(http.StatusOK, response)
}