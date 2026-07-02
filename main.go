package main

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/middleware"
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

	// Connect to the database
	config.Connect()

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	router.GET("/api/test-success", func(c *gin.Context) {
		utils.SendSuccess(c, http.StatusOK, "Berhasil mengambil data uji", gin.H{"id": 1, "name": "PPDB Test"})
	})

	router.GET("/api/test-error", func(c *gin.Context) {
		utils.SendError(c, http.StatusBadRequest, utils.MsgBadRequest, utils.ErrInvalidInput)
	})

	// Run the server on port 8080
	err = router.Run(":8080")
	if err != nil {
		log.Println("Gagal menjalankan server:", err)
	}
}
