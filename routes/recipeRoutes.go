package routes

import (
	"book-recipe-be-go/controllers"

	"github.com/gin-gonic/gin"
)

func RecipeRoutes(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/book-recipe/book-recipes", controllers.GetAllRecipes)
	apiGroup.PUT("/book-recipe/book-recipes/:recipeId/favorites", controllers.ToggleFavorite)
	apiGroup.GET("/book-recipe/book-recipes/:recipeId", controllers.GetRecipeDetailsById)
}