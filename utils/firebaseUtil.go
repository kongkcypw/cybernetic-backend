package utils

import (
	database "example/backend/database"
	"fmt"
	"log"
	"os"
)

// DownloadImageFile downloads an image file from Firebase Storage
func DownloadImageFile(filePath, destPath string) error {
	// Example: Get a file from Firebase Storage
	bucketName := os.Getenv("FIREBASE_STORAGE_BUCKET")
	if err := database.GetFileFromBucket(bucketName, filePath, destPath); err != nil {
		log.Fatalf("failed to get file from bucket: %v", err)
		return err
	}
	return nil
}

// UploadImageFile uploads an image file to Firebase Storage
func UploadImageFile(srcPath, destPath string) error {
	// Example: Upload a file to Firebase Storage
	bucketName := os.Getenv("FIREBASE_STORAGE_BUCKET")
	if err := database.UploadFileToBucket(bucketName, srcPath, destPath); err != nil {
		log.Fatalf("failed to upload file to bucket: %v", err)
		return err
	}
	return nil
}

func GetStorageFileURL(filePath string) (string, error) {
	bucketName := os.Getenv("FIREBASE_STORAGE_BUCKET")
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filePath)
	return url, nil
}
