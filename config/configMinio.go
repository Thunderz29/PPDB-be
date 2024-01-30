package config

import (
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

