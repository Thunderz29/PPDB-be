package routes

import (
	"ppdb-be/controllers"
	"ppdb-be/middleware"

	"ppdb-be/repositories"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, userCtrl controllers.UserController, sessionRepo repositories.SessionRepository) {
	router.POST("/login", userCtrl.Login)
	router.POST("/register", userCtrl.CreateUser)

	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(sessionRepo))
	{
		protected.POST("/logout", userCtrl.Logout)

		users := protected.Group("/users")
		{
			users.GET("", userCtrl.GetUsers)
			users.GET("/:id", userCtrl.GetUserByID)
			users.PUT("/:id", userCtrl.UpdateUser)
			users.DELETE("/:id", userCtrl.DeleteUser)
		}
	}
}
