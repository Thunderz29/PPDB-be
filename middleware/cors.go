package middleware

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func SetupCORS(router *gin.Engine) {
    // Configure CORS middleware
    middle := cors.DefaultConfig()
    middle.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000", "https://mt-react-golang.cloudias79.com"}
    middle.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    middle.AllowHeaders = []string{"Authorization", "Content-Type"}

    // Use CORS middleware with the configured options
    router.Use(cors.New(middle))
}
