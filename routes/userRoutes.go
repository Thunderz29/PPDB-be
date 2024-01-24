package routes

import (
	"book-recipe-be-go/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/api/user-management/users/sign-up",controllers.SignUp)
	router.POST("/api/user-management/users/signin",controllers.SignIn)
}