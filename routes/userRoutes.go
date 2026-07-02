package routes

import (
	"ppdb-be/controllers"
	"ppdb-be/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.CreateUser)

	protected := router.Group("/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", controllers.GetUsers)
		protected.GET("/:id", controllers.GetUserByID)
		protected.PUT("/:id", controllers.UpdateUser)
		protected.DELETE("/:id", controllers.DeleteUser)
	}
}
