package routes

import (
    "book-recipe-be-go/controllers"
    "github.com/gin-gonic/gin"
)

func OptionListRoutes(apiGroup *gin.RouterGroup) {
    apiGroup.GET("/book-recipe-masters/category-option-lists", controllers.GetAllCategories)
    apiGroup.GET("/book-recipe-masters/level-option-lists", controllers.GetAllLevels)
}
