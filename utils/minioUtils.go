package utils

import (
	"book-recipe-be-go/config"
	"book-recipe-be-go/models/request"
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"github.com/minio/minio-go/v7"
)

func UploadFileToMinio(file *multipart.FileHeader, recipeRequest interface{}) (string, error) {
	minioClient, err := config.ConfigMinio()
	if err != nil {
		return "", fmt.Errorf("Failed to initialize MinIO client: %v", err)
	}

	ctx := context.Background()
	bucketName := "talent79-dev"

	var recipeName, categoryName, levelName string
	var timestamp, fileExtension, generatedFilename string

	switch req := recipeRequest.(type) {
	case *request.UpdateRecipeRequest:
		// Use fields from UpdateRecipeRequest
		recipeName = sanitizeForFilename(req.RecipeName)
		categoryName = sanitizeForFilename(req.Categories.CategoryName)
		levelName = sanitizeForFilename(req.Levels.LevelName)
		timestamp = strconv.FormatInt(time.Now().UnixNano(), 10)
		fileExtension = getFileExtension(file.Filename)
		generatedFilename = fmt.Sprintf("%s_%s_%s_%s%s", recipeName, categoryName, levelName, timestamp, fileExtension)
	case *request.CreateRecipeRequest:
		// Use fields from CreateRecipeRequest
		recipeName = sanitizeForFilename(req.RecipeName)
		categoryName = sanitizeForFilename(req.Categories.CategoryName)
		levelName = sanitizeForFilename(req.Levels.LevelName)
		timestamp = strconv.FormatInt(time.Now().UnixNano(), 10)
		fileExtension = getFileExtension(file.Filename)
		generatedFilename = fmt.Sprintf("%s_%s_%s_%s%s", recipeName, categoryName, levelName, timestamp, fileExtension)
	default:
		return "", fmt.Errorf("Unsupported request type")
	}

	fileData, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("Failed to open file: %v", err)
	}
	defer fileData.Close()

	_, err = minioClient.PutObject(ctx, bucketName, generatedFilename, fileData, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("Failed to upload file to MinIO: %v", err)
	}

	return generatedFilename, nil
}


func getMinioURL(minioClient *minio.Client, filename string) (string, error) {
	// Modify this function to generate Minio URL using minioClient
	bucketName := "talent79-dev" // Change this to your bucket name
	const DEFAULT_EXPIRY = 3600

	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, filename, time.Second*DEFAULT_EXPIRY, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func GetPublicLink(minioClient *minio.Client, filename string) (string, error) {
	return getMinioURL(minioClient, filename)
}


// sanitizeForFilename membersihkan string untuk digunakan sebagai nama file
func sanitizeForFilename(input string) string {
    // Ganti karakter yang tidak valid dengan underscore
    return strings.Map(func(r rune) rune {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            return r
        }
        return '_'
    }, input)
}

// getFileExtension mendapatkan ekstensi file dari nama file
func getFileExtension(filename string) string {
    dotIndex := strings.LastIndex(filename, ".")
    if dotIndex == -1 {
        return ""
    }
    return filename[dotIndex:]
}