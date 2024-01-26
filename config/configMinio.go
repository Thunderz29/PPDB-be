package config

import (
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
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConfigMinio() (*minio.Client, error) {
	endpoint := "minio.cloudias79.com"
	accessKeyID := "talent79"
	secretAccessKey := "evT25hDUUdf0gl1M9WOtwo3T3"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
        return nil, err
    }

	return minioClient, nil
}

func UploadFileToMinio(file *multipart.FileHeader, createRecipeRequest *request.CreateRecipeRequest) (string, error) {
    minioClient, err := ConfigMinio()
    if err != nil {
        return "", fmt.Errorf("Failed to initialize MinIO client: %v", err)
    }

    ctx := context.Background()
    bucketName := "talent79-dev"

    // Cleanse strings for filename
    recipeName := sanitizeForFilename(createRecipeRequest.RecipeName)
    categoryName := sanitizeForFilename(createRecipeRequest.Categories.CategoryName)
    levelName := sanitizeForFilename(createRecipeRequest.Levels.LevelName)

    if recipeName == "" || categoryName == "" || levelName == "" {
        return "", fmt.Errorf("One or more components for filename are empty. Recipe: %s, Category: %s, Level: %s",
            createRecipeRequest.RecipeName, createRecipeRequest.Categories.CategoryName, createRecipeRequest.Levels.LevelName)
    }

    timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
    fileExtension := getFileExtension(file.Filename)

    generatedFilename := fmt.Sprintf("%s_%s_%s_%s%s", recipeName, categoryName, levelName, timestamp, fileExtension)

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