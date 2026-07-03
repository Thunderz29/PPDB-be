package main

import (
	"ppdb-be/config"
	_ "ppdb-be/docs"
	"ppdb-be/middleware"
	"ppdb-be/routes"
	"log"
	"net/http"

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
func main() {
	logFile, err := config.InitLogger()
	if err != nil {
		log.Println("Gagal menginisialisasi logger:", err)
	}
	defer logFile.Close()

	router := gin.New()

	middleware.SetupCORS(router)

	config.Connect()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	apiGroup := router.Group("/api")
	{
		routes.UserRoutes(apiGroup)
	}

	err = router.Run(":8080")
	if err != nil {
		log.Println("Gagal menjalankan server:", err)
	}
}
