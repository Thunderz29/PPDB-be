package controllers

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"book-recipe-be-go/models/request"
	"book-recipe-be-go/models/response"
	"book-recipe-be-go/utils"
	"encoding/json"

	// "fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func DeleteMyRecipe(c *gin.Context) {
	recipeID := c.Param("recipeId")
	userID := c.Query("userId")
	var existingRecipe models.Recipe

	recipeIDInt, err := strconv.Atoi(recipeID)
	if err != nil {
		response := response.MyRecipeResponse{
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
			StatusCode: http.StatusInternalServerError,
			Details:    err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		response := response.MyRecipeResponse{
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
			StatusCode: http.StatusInternalServerError,
			Details:    err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	recipeName, err := utils.GetRecipeNameByRecipeID(recipeIDInt)
	if err != nil {
		response := response.DataResponse{
			Total:      0,
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	err = config.DB.Where("recipe_id = ? AND user_id = ? AND is_deleted = ?", recipeIDInt, userIDInt, false).First(&existingRecipe).Error
	if err != nil {
		response := response.MyRecipeResponse{
			Message:    "Data Already Deleted!",
			StatusCode: http.StatusBadRequest,
			Details: recipeName + " sudah terhapus!",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = config.DB.Model(&models.Recipe{}).
		Where("recipe_id = ? AND user_id = ?", recipeIDInt, userIDInt).Update("is_deleted", true).Error
	if err != nil {
		response := response.MyRecipeResponse{
			Message:    "Terjadi kesalahan saat menghapus data.",
			StatusCode: http.StatusInternalServerError,
			Details:    err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.RecipeDetailsResponse{
		Total: 1,
		Data: nil,
		Message:    "Resep " + recipeName + " berhasil dihapus!",
		StatusCode: http.StatusOK,
		Status: "OK",
	}
	c.JSON(http.StatusOK, response)

}

func GetAllMyRecipes(c *gin.Context) {
	var recipes []models.Recipe
	var recipeFilter request.RecipeFilter
	var total int64

	minioClient, err := config.ConfigMinio()
	if err != nil {
		log.Fatalln(err)
		return
	}

	if err := c.ShouldBindQuery(&recipeFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengonversi nilai string ke integer
	pageNumberInt, _ := strconv.Atoi(recipeFilter.PageNumber)
	pageSizeInt, _ := strconv.Atoi(recipeFilter.PageSize)

	// Menetapkan nilai default jika tidak ada nilai atau nilai tidak valid
	if pageNumberInt <= 0 {
		pageNumberInt = 1
	}

	if pageSizeInt <= 0 {
		pageSizeInt = 10 // Atur ke ukuran halaman default
	}

	if recipeFilter.UserID == "" {
		response := response.MyRecipeResponse{
			Message:    "User tidak valid",
			StatusCode: http.StatusUnauthorized,
			Details:    "Unauthorized",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db := config.DB.Model(&models.Recipe{})

	// Apply filters
	db = utils.ApplyRecipeFilters(db, &recipeFilter)

	// Filter by time cook
	if recipeFilter.Time != "" {
		timeCookInt, err := strconv.Atoi(recipeFilter.Time)

		if err != nil {
			response := response.MyRecipeResponse{
				Message:    "Parameter time harus berupa angka",
				StatusCode: http.StatusBadRequest,
				Details:    "Bad Request",
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if timeCookInt <= 30 {
			db = db.Where("time_cook <= ?", timeCookInt)
		} else if timeCookInt > 30 && timeCookInt <= 60 {
			db = db.Where("time_cook > 30 AND time_cook <= ?", timeCookInt)
		} else if timeCookInt > 60 {
			db = db.Where("time_cook > 60 AND time_cook <= ?", timeCookInt)
		}

	}

	userIDInt, err := strconv.Atoi(recipeFilter.UserID)
	if err != nil {
		response := response.MyRecipeResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Details:    "Bad Request",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Menggunakan Joins untuk mengambil data dari tabel levels dan categories
	if err := db.
		Preload("Category").
		Preload("Level").
		Where("user_id = ?", userIDInt).
		Order("LOWER(recipe_name) ASC").
		Where("is_deleted = ?", false).
		Count(&total).
		Limit(pageSizeInt).Offset((pageNumberInt - 1) * pageSizeInt).
		Find(&recipes).Error; err != nil {
		response := response.MessageResponse{
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
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
		response := response.MyRecipeResponse{
			Message:    "Resep masakan tidak tersedia",
			StatusCode: http.StatusNotFound,
			Details:    "Data tidak ditemukan!",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := response.DataResponse{
		Total:      total,
		Data:       recipeEntries,
		Message:    "Berhasil memuat Resep Masakan Saya",
		StatusCode: http.StatusOK,
		Status:     "OK",
	}
	c.JSON(http.StatusOK, response)
}

func GetAllRecipes(c *gin.Context) {
	var recipes []models.Recipe
	var recipeFilter request.RecipeFilter
	var total int64

	minioClient, err := config.ConfigMinio()
	if err != nil {
		log.Fatalln(err)
		return
	}

	if err := c.ShouldBindQuery(&recipeFilter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengonversi nilai string ke integer
	pageNumberInt, _ := strconv.Atoi(recipeFilter.PageNumber)
	pageSizeInt, _ := strconv.Atoi(recipeFilter.PageSize)

	// Menetapkan nilai default jika tidak ada nilai atau nilai tidak valid
	if pageNumberInt <= 0 {
		pageNumberInt = 1
	}

	if pageSizeInt <= 0 {
		pageSizeInt = 10 // Atur ke ukuran halaman default
	}

	if recipeFilter.UserID == "" {
		response := response.MessageResponse{
			Message:    "User tidak valid",
			StatusCode: http.StatusUnauthorized,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db := config.DB.Model(&models.Recipe{})

	// Apply filters
	db = utils.ApplyRecipeFilters(db, &recipeFilter)

	// Filter by time cook
	if recipeFilter.Time != "" {
		if recipeFilter.Time == "60" {
			// Jika waktu adalah '60', cari resep dengan waktu memasak tepat 60
			db = db.Where("time_cook > ?", 60)
		} else {
			// Jika waktu tidak sama dengan '60', lakukan pemrosesan seperti sebelumnya
			timeCookRange := strings.Split(recipeFilter.Time, "-")
			if len(timeCookRange) == 2 {
				minTimeCook, err := strconv.Atoi(timeCookRange[0])
				maxTimeCook, err := strconv.Atoi(timeCookRange[1])

				if err != nil {
					response := response.MessageResponse{
						Message:    err.Error(),
						StatusCode: http.StatusBadRequest,
						Status:     "ERROR",
					}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				if minTimeCook >= 0 && maxTimeCook > minTimeCook {
					db = db.Where("time_cook >= ? AND time_cook <= ?", minTimeCook, maxTimeCook)
				} else {
					response := response.MessageResponse{
						Message:    "Invalid time cook range",
						StatusCode: http.StatusBadRequest,
						Status:     "ERROR",
					}
					c.JSON(http.StatusBadRequest, response)
					return
				}
			} else {
				response := response.MessageResponse{
					Message:    "Invalid time cook format",
					StatusCode: http.StatusBadRequest,
					Status:     "ERROR",
				}
				c.JSON(http.StatusBadRequest, response)
				return
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

		userIDInt, err := strconv.Atoi(recipeFilter.UserID)
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
		Message:    "Berhasil memuat Resep Masakan",
		StatusCode: http.StatusOK,
		Status:     "Success",
	}
	c.JSON(http.StatusOK, response)
}

func CreateRecipe(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response := response.MessageResponse{
			Message:    "Error reading file from form-data",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Baca data dari form-data
	var request request.CreateRecipeRequest
	jsonStr := c.PostForm("request")
	if err := json.Unmarshal([]byte(jsonStr), &request); err != nil {
		response := response.MessageResponse{
			Message:    "Error parsing JSON from form-data",
			StatusCode: http.StatusBadRequest,
			Status:     "ERROR",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Mapping data ke model Recipe
	categories := models.Category{
		CategoryID:   request.Categories.CategoryId,
		CategoryName: request.Categories.CategoryName,
	}

	levels := models.Level{
		LevelID:   request.Levels.LevelId,
		LevelName: request.Levels.LevelName,
	}

	username, err := utils.GetusernameByUserID(uint(request.UserId))
	if err != nil {
		response := response.MessageResponse{
			Message:    "Error getting user information",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	imageFilename, err := config.UploadFileToMinio(file, &request)
	if err != nil {
		response := response.MessageResponse{
			Message:    "Error uploading image to Minio",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	recipe := models.Recipe{
		CategoryID:    request.Categories.CategoryId,
		UserID:        request.UserId,
		LevelID:       request.Levels.LevelId,
		RecipeName:    request.RecipeName,
		ImageFilename: imageFilename,
		TimeCook:      &request.TimeCook,
		Ingridient:    request.Ingridient,
		HowToCook:     request.HowToCook,
		IsDeleted:     false,
		CreatedBy:     username,
		CreatedTime:   time.Now(),
		ModifiedBy:    username,
		ModifiedTime:  time.Now(),
		Category:      categories,
		Level:         levels,
	}

	if err := config.DB.Create(&recipe).Error; err != nil {
		response := response.MessageResponse{
			Message:    "Terjadi kesalahan server. Silakan coba kembali",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Respons berhasil
	responseMessage := "Resep " + request.RecipeName + " berhasil ditambahkan!"
	c.JSON(http.StatusOK, response.MessageResponse{
		Message:    responseMessage,
		StatusCode: http.StatusOK,
		Status:     "OK",
	})
}

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
			Total:      0,
			Data:       nil,
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
			Total:      0,
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
				Total:      0,
				Data:       nil,
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
			Total:      1,
			Data:       nil,
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
			Total:      0,
			Data:       nil,
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := response.DataResponse{
		Total:      1,
		Data:       nil,
		Message:    "Resep " + recipeName + " berhasil ditambahkan ke dalam favorit",
		StatusCode: http.StatusOK,
		Status:     "OK",
	}
	c.JSON(http.StatusOK, response)
}

func GetRecipeDetailsById(c *gin.Context) {
	recipeID := c.Param("recipeId")
	var recipeDetails models.Recipe
	var total int64

	minioClient, err := config.ConfigMinio()
	if err != nil {
		log.Fatalln(err)
		return
	}

	recipeIDInt, _ := strconv.Atoi(recipeID)

	// Mendapatkan token dari header HTTP
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		response := response.DataResponse{
			Total:      0,
			Data:       nil,
			Message:    "Token tidak ditemukan",
			StatusCode: http.StatusUnauthorized,
			Status:     "Unauthorized",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Men-decode token untuk mendapatkan informasi pengguna
	userId, err := utils.GetUserIdFromToken(tokenString)
	if err != nil {
		response := response.DataResponse{
			Total:      0,
			Data:       nil,
			Message:    "Token tidak valid",
			StatusCode: http.StatusUnauthorized,
			Status:     "Unauthorized",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Menggunakan Joins untuk mengambil data dari tabel levels dan categories
	db := config.DB.Model(&models.Recipe{})
	if err := db.
		Preload("Category").
		Preload("Level").
		Where("recipe_id = ?", recipeIDInt).
		Count(&total).
		Find(&recipeDetails).Error; err != nil {
		response := response.DataResponse{
			Total:      0,
			Data:       nil,
			Message:    "Terjadi kesalahan server. Silakan coba kembali.",
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if total == 0 {
		response := response.RecipeDetailsResponse{
			Total:      0,
			Data:       nil,
			Message:    "Detil Resep masakan tidak tersedia",
			StatusCode: http.StatusNotFound,
			Status:     "Not Found",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	categoryName, err := utils.GetCategoryName(recipeDetails.Category.CategoryID)
	if err != nil {
		log.Println("Error fetching category name:", err)
		categoryName = "Unknown Category"
	}

	levelName, err := utils.GetLevelName(recipeDetails.Level.LevelID)
	if err != nil {
		log.Println("Error fetching level name:", err)
		levelName = "Unknown Level"
	}

	imageUrl, err := utils.GetImageURL(minioClient, recipeDetails.ImageFilename)
	if err != nil {
		log.Println("Error getting image URL:", err)
	}

	isFavorite := utils.CheckFavoriteRecipe(uint(userId), recipeIDInt)

	// Modify the response mapping
	data := response.RecipeDetailsEntry{
		RecipeId:      recipeDetails.RecipeID,
		Categories:    response.CategoryInfo{CategoryId: recipeDetails.CategoryID, CategoryName: categoryName},
		Levels:        response.LevelInfo{LevelId: recipeDetails.LevelID, LevelName: levelName},
		RecipeName:    recipeDetails.RecipeName,
		ImageFilename: imageUrl,
		TimeCook:      recipeDetails.TimeCook,
		Ingridient:    recipeDetails.Ingridient,
		HowToCook:     recipeDetails.HowToCook,
		IsFavorite:    isFavorite,
	}

	response := response.RecipeDetailsResponse{
		Total:      total,
		Data:       data,
		Message:    "Berhasil memuat Resep Masakan",
		StatusCode: http.StatusOK,
		Status:     "Success",
	}
	c.JSON(http.StatusOK, response)
}
