package utils

import (
    "book-recipe-be-go/config"
    "book-recipe-be-go/models"
)

func CheckFavoriteRecipe(userID uint, recipeID int) bool {
    var favoriteFood models.FavoriteFood

    if err := config.DB.Where("user_id = ? AND recipe_id = ?", userID, recipeID).First(&favoriteFood).Error; err != nil {
        // Jika terdapat error, diasumsikan bahwa resep bukan favorit
        return false
    }

    // Jika tidak ada error, mengembalikan nilai is_favorite pada tabel favorite_foods
    return favoriteFood.IsFavorite
}