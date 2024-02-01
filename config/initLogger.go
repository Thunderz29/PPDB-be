package config

import (
	"log"
	"os"
)

func InitLogger() (*os.File, error) {
	// Buat atau buka file log
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Gagal membuka file log:", err)
	}

	// Set output log ke file
	log.SetOutput(file)

	return file, err
}