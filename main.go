package main

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	config.Connect()
	routes.UserRoutes(router)
	router.Run(":8080")
}