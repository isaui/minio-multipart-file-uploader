package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient is the global MinIO client instance
var MinioClient *minio.Client

// MinioCore is the global MinIO core client instance for multipart operations
var MinioCore *minio.Core

// InitMinio initializes the MinIO client connection
func InitMinio() (*minio.Client, error) {
	// Get MinIO configuration from environment variables
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}
	
	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Printf("Failed to create MinIO client: %v", err)
		return nil, err
	}
	
	// Ensure the bucket exists
	bucketName := "uploads"
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("Error checking if bucket exists: %v", err)
		return nil, err
	}
	
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Error creating bucket: %v", err)
			return nil, fmt.Errorf("error creating bucket: %v", err)
		}
		log.Printf("Successfully created bucket: %s", bucketName)
	} else {
		log.Printf("Bucket '%s' already exists", bucketName)
	}
	
	MinioClient = client
	
	// Initialize the core client for multipart operations
	MinioCore = &minio.Core{Client: client}
	
	return client, nil
}
