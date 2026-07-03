package routes

import (
	"ppdb-be/controllers"
	"ppdb-be/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.CreateUser)

	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/logout", controllers.Logout)

		users := protected.Group("/users")
		{
			users.GET("", controllers.GetUsers)
			users.GET("/:id", controllers.GetUserByID)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
		}
	}
}
