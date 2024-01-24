package routes

import (
	"book-recipe-be-go/controllers"

	"github.com/gin-gonic/gin"
)

func RecipeRoutes(router *gin.Engine) {
	router.GET("/book-recipe/book-recipes", controllers.GetAllRecipes)
	router.PUT("/book-recipe/book-recipes/:recipeId/favorites", controllers.ToggleFavorite)
}