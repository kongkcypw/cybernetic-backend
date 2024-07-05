package controllers

import (
	"example/backend/database"

	"github.com/gin-gonic/gin"
)

func UploadImageToFirebase() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("image")
		path := c.Query("path")
		if err != nil {
			c.JSON(400, gin.H{
				"message": "no file found",
			})
			return
		}
		err = database.UploadImageToFirebase(file, header.Filename, path)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "upload image failed",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "upload image success",
		})
	}
}
