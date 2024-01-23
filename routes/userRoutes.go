package routes

import (
	"book-recipe-be-go/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/user-management/users/sign-up",controllers.SignUp)
}