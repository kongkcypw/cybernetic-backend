package controllers

import (
	"context"
	"log"

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

		var character models.UserCharacter
		err := collection.FindOne(context.Background(), bson.M{"userId": userId}).Decode(&character)
		if err != nil {
			c.JSON(404, gin.H{"error": "Character not found!"})
			return
		}

		c.JSON(200, gin.H{"character": character})
	}
}

func GetUserGamePlayLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("userId")
		log.Println("userId: ", userId)

		if userId == "" {
			c.JSON(400, gin.H{"error": "userId is required"})
			return
		}

		client := database.MongoDB()
		if client == nil {
			c.JSON(500, gin.H{"error": "Database connection is not initialized"})
			return
		}

		collection := database.MongoDBOpenCollection(client, "cybernetic", "user_game_play_level")

		var userGamePlayLevel models.UserGamePlayLevel
		err := collection.FindOne(context.Background(), bson.M{"userId": userId}).Decode(&userGamePlayLevel)
		if err != nil {
			c.JSON(404, gin.H{"error": "Character not found!", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{"userGamePlayLevel": userGamePlayLevel})
	}
}

func UpdateLevelPlayed() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			UserId      string `json:"userId"`
			LevelNumber int64  `json:"levelNumber"` // Use int64 for level numbers
		}

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		userId := requestBody.UserId
		levelPlayed := requestBody.LevelNumber

		if userId == "" {
			c.JSON(400, gin.H{"error": "userId is required"})
			return
		}

		client := database.MongoDB()
		if client == nil {
			c.JSON(500, gin.H{"error": "Database connection is not initialized"})
			return
		}

		collection := database.MongoDBOpenCollection(client, "cybernetic", "user_game_play_level")

		filter := bson.M{"userId": userId}
		update := bson.M{"$addToSet": bson.M{"levelPlayed": levelPlayed}} // Use $addToSet to ensure uniqueness

		_, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update levelPlayed", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "LevelPlayed updated successfully"})
	}
}
