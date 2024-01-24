// main.go
package main

import (
    "book-recipe-be-go/config"
    "book-recipe-be-go/middleware"
    "book-recipe-be-go/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    // Create a new Gin router
    router := gin.New()

    // Use CORS middleware
    middleware.SetupCORS(router)

    // Connect to the database (assuming Connect function is defined in the config package)
    config.Connect()

    // Define user-related routes using the UserRoutes function from the routes package
    routes.UserRoutes(router)

    // Run the server on port 8080
    router.Run(":8080")
}
