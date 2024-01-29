package routes

import (
	"book-recipe-be-go/controllers"

	"github.com/gin-gonic/gin"
)

func RecipeRoutes(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/book-recipe/book-recipes", controllers.GetAllRecipes)
	apiGroup.POST("/book-recipe/book-recipes", controllers.CreateRecipe)
	apiGroup.PUT("/book-recipe/book-recipes", controllers.UpdateRecipe)
	apiGroup.PUT("/book-recipe/book-recipes/:recipeId/favorites", controllers.ToggleFavorite)
	apiGroup.GET("/book-recipe/book-recipes/:recipeId", controllers.GetRecipeDetailsById)
	apiGroup.GET("/book-recipe/my-recipes", controllers.GetAllMyRecipes)
	apiGroup.PUT("/book-recipe/book-recipes/:recipeId", controllers.DeleteMyRecipe)
	apiGroup.GET("/book-recipe/my-favorite-recipes", controllers.GetAllMyFavRecipes)
}