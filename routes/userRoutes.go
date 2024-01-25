package routes

import (
    "book-recipe-be-go/controllers"
    "github.com/gin-gonic/gin"
)

func UserRoutes(apiGroup *gin.RouterGroup) {
    apiGroup.POST("/user-management/users/sign-up", controllers.SignUp)
    apiGroup.POST("/user-management/users/signin", controllers.SignIn)
}
