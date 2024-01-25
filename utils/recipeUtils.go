package utils

import (
    "book-recipe-be-go/config"
    "book-recipe-be-go/models"
	"github.com/minio/minio-go/v7"
)

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