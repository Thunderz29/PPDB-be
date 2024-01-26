package utils

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models"
	"book-recipe-be-go/models/request"
	"fmt"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func ApplyRecipeFilters(db *gorm.DB, filter *request.RecipeFilter) *gorm.DB {
	// Filter by recipe name
	if filter.RecipeName != "" {
		db = db.Where("recipe_name LIKE ?", "%"+filter.RecipeName+"%")
	}

	// Filter by level id
	if filter.LevelID != "" {
		levelID, err := strconv.Atoi(filter.LevelID)
		if err == nil {
			db = db.Where("level_id = ?", levelID)
		}
	}

	// Filter by category id
	if filter.CategoryID != "" {
		categoryID, err := strconv.Atoi(filter.CategoryID)
		if err == nil {
			db = db.Where("category_id = ?", categoryID)
		}
	}

	// Menentukan pengurutan berdasarkan sortBy
	if filter.SortBy != "" {
		sortColumns := map[string]string{
			"recipeName": "LOWER(recipe_name)",
			"timeCook":   "time_cook",
		}

		sortParams := strings.Split(filter.SortBy, ",")
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

	return db
}

func GetCategoryName(categoryID int) (string, error) {
	var category models.Category
	if err := config.DB.Model(&models.Category{}).Where("category_id = ?", categoryID).First(&category).Error; err != nil {
		return "", err
	}
	return category.CategoryName, nil
}

func GetLevelName(levelID int) (string, error) {
	var level models.Level
	if err := config.DB.Model(&models.Level{}).Where("level_id = ?", levelID).First(&level).Error; err != nil {
		return "", err
	}
	return level.LevelName, nil
}

func GetImageURL(minioClient *minio.Client, filename string) (string, error) {
	url, err := config.GetPublicLink(minioClient, filename)
	if err != nil {
		return "", err
	}

	return url, nil
}

func GetRecipeNameByRecipeID(recipeID int) (string, error) {
	var recipe models.Recipe
	if err := config.DB.Model(&models.Recipe{}).Where("recipe_id = ?", recipeID).First(&recipe).Error; err != nil {
		return "", err
	}
	return recipe.RecipeName, nil
}