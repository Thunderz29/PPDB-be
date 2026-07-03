package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS(router *gin.Engine) {
	// Configure CORS middleware
	middle := cors.DefaultConfig()

	originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if originsEnv != "" {
		middle.AllowOrigins = strings.Split(originsEnv, ",")
	} else {
		// Fallback for development if env not set
		middle.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	}

	middle.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	middle.AllowHeaders = []string{"Authorization", "Content-Type"}

	// Use CORS middleware with the configured options
	router.Use(cors.New(middle))
}
