package main

import (
	"log/slog"
	"net/http"
	"ppdb-be/config"
	"ppdb-be/controllers"
	_ "ppdb-be/docs"
	"ppdb-be/middleware"
	"ppdb-be/repositories"
	"ppdb-be/routes"
	"ppdb-be/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title PPDB API
// @version 1.0
// @description API Server for PPDB Application
// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// HealthCheck godoc
// @Summary Check server status
// @Description Endpoint to check if the server is running and ready to accept requests.
// @Tags System
// @Produce json
// @Success 200 {object} map[string]interface{} "Status OK"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Server PPDB Berjalan Lancar"})
}
func main() {
	config.InitLogger()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())

	middleware.SetupCORS(router)

	config.Connect()

	userRepo := repositories.NewUserRepository(config.DB)
	sessionRepo := repositories.NewSessionRepository(config.DB)
	userService := services.NewUserService(userRepo, sessionRepo)
	userController := controllers.NewUserController(userService)

	roleRepo := repositories.NewRoleRepository(config.DB)
	roleService := services.NewRoleService(roleRepo)
	roleController := controllers.NewRoleController(roleService)

	profileRepo := repositories.NewSchoolProfileRepository(config.DB)
	profileService := services.NewSchoolProfileService(profileRepo)
	profileController := controllers.NewSchoolProfileController(profileService)

	configRepo := repositories.NewConfigRepository(config.DB)
	configService := services.NewConfigService(configRepo)
	configController := controllers.NewConfigController(configService)

	schoolRepo := repositories.NewMasterSchoolRepository(config.DB)
	schoolService := services.NewMasterSchoolService(schoolRepo)
	schoolController := controllers.NewMasterSchoolController(schoolService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api/health", HealthCheck)

	apiGroup := router.Group("/api")
	{
		routes.UserRoutes(apiGroup, userController, sessionRepo)
		routes.MasterRoutes(apiGroup, roleController, profileController, configController, schoolController, sessionRepo)
	}

	err := router.Run(":8080")
	if err != nil {
		slog.Error("Gagal menjalankan server:", "error", err)
	}
}
