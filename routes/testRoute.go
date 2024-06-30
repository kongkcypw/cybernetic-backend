package routes

import (
	utils "example/backend/utils"

	"github.com/gin-gonic/gin"
)

func TestRoute(router *gin.Engine) {
	router.GET("/test/firebase/download-image-file", func(c *gin.Context) {
		err := utils.DownloadImageFile("orangecat.jpg", "orangecat.jpg")
		if err != nil {
			c.JSON(500, gin.H{
				"message": "download image failed",
			})
		}
		c.JSON(200, gin.H{
			"message": "donwload image success",
		})
	})
	router.GET("/test/firebase/upload-image-file", func(c *gin.Context) {
		err := utils.UploadImageFile("orangecat.jpg", "test/orangecat.jpg")
		if err != nil {
			c.JSON(500, gin.H{
				"message": "upload image failed",
			})
		}
		c.JSON(200, gin.H{
			"message": "upload image success",
		})
	})
	router.GET("/test/firebase/get-storage-file-url", func(c *gin.Context) {
		// url, err := utils.GetStorageFileURL("test/orangecat.jpg")
		url, err := utils.GetStorageFileURL("scene_environments/items/scene_name/Tank.gltf")
		if err != nil {
			c.JSON(500, gin.H{
				"message": "get file URL failed",
			})
		}
		c.JSON(200, gin.H{
			"message": "get file URL success",
			"url":     url,
		})
	})
}
