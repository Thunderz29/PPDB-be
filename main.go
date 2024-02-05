package main

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/middleware"
	"book-recipe-be-go/routes"
	"book-recipe-be-go/utils"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
    logFile, err := config.InitLogger()
	if err != nil {
		log.Println("Gagal menginisialisasi logger:", err)
	}
	defer logFile.Close()

    
    // Create a new Gin router
    router := gin.New()

    // Use CORS middleware
    middleware.SetupCORS(router)

    // Connect to the database (assuming Connect function is defined in the config package)
    config.Connect()

    log.Println(utils.WelcomeMessagee)

    router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status":"OK"})
	})

    apiGroup := router.Group("/api")
    {
        routes.UserRoutes(apiGroup)
        routes.RecipeRoutes(apiGroup)
        routes.OptionListRoutes(apiGroup)
    }

    // Run the server on port 8080
	err = router.Run(":8080")
	if err != nil {
		log.Println(utils.ErrInternalServer, err)
	}

}
