package routes

import (
	database "example/backend/database"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func TestRoute(router *gin.Engine) {
	router.GET("/test/firebase/download-image-file", func(c *gin.Context) {
		// Example: Get a file from Firebase Storage
		bucketName := os.Getenv("FIREBASE_STORAGE_BUCKET")
		filePath := "orangecat.jpg"
		destPath := "orangecat.jpg"
		if err := database.GetFileFromBucket(bucketName, filePath, destPath); err != nil {
			log.Fatalf("failed to get file from bucket: %v", err)
		}
		c.JSON(200, gin.H{
			"message": "donwload image success",
		})
	})
}
