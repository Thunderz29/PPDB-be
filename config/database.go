package config

import (
	// "book-recipe-be-go/models"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
        panic(err)
    }

    // Read values from environment variables
    host := os.Getenv("DB_HOST")
    portStr := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")

    // Parse port string to integer
    port, err := strconv.Atoi(portStr)
    if err != nil {
		panic(err)
    }

    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }

	// db.AutoMigrate(&models.User{})
	// db.AutoMigrate(&models.Category{})
	// db.AutoMigrate(&models.Level{})
	// db.AutoMigrate(&models.Recipe{})
	// db.AutoMigrate(&models.FavoriteFood{})
	DB = db
}

