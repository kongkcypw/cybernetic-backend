package controllers

import (

	// "log"
	// "strconv"
	// "net/http"
	// "time"
	// "fmt"

	"github.com/gin-gonic/gin"

	orm "example/backend/config"
	"example/backend/models"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("userId")

		if userId == "" {
			c.JSON(400, gin.H{"error": "userId is required"})
			return
		}

		var user models.User
		dbErr := orm.DB().Where("userId = ?", userId).First(&user).Error

		if dbErr != nil {
			c.JSON(404, gin.H{"error": "User not found!"})
			return
		}
		c.JSON(200, user)
	}
}
