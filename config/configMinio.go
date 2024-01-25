package config

import (
	"context"
	"net/url"
	"time"

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