package controllers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	database "example/backend/database"
	"example/backend/models"
)

func GetCharacter() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("userId")

		if userId == "" {
			c.JSON(400, gin.H{"error": "userId is required"})
			return
		}

		client := database.MongoDB()
		if client == nil {
			c.JSON(500, gin.H{"error": "Database connection is not initialized"})
			return
		}

		collection := database.MongoDBOpenCollection(client, "cybernetic", "user_character")

		var character models.Character
		err := collection.FindOne(context.Background(), bson.M{"userId": userId}).Decode(&character)
		if err != nil {
			c.JSON(404, gin.H{"error": "Character not found!"})
			return
		}

		c.JSON(200, gin.H{"character": character})
	}
}
