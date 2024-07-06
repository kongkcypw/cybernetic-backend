package controllers

import (
	"context"
	database "example/backend/database"
	"example/backend/models"
	"log"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetGamePlayLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := database.MongoDB()
		if client == nil {
			c.JSON(500, gin.H{"error": "Database connection is not initialized"})
			return
		}
		collection := database.MongoDBOpenCollection(client, "cybernetic", "game_play_level")
		var gamePlayLevel []models.GamePlayLevel
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching game play levels"})
			return
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var level models.GamePlayLevel
			cursor.Decode(&level)
			gamePlayLevel = append(gamePlayLevel, level)
		}
		// Sort gamePlayLevel by level_number
		sort.Slice(gamePlayLevel, func(i, j int) bool {
			return gamePlayLevel[i].LevelNumber < gamePlayLevel[j].LevelNumber
		})

		c.JSON(200, gin.H{"gamePlayLevel": gamePlayLevel})
	}
}

func GetGamePlayLevelByNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := database.MongoDB()
		if client == nil {
			c.JSON(500, gin.H{"error": "Database connection is not initialized"})
			return
		}

		collection := database.MongoDBOpenCollection(client, "cybernetic", "game_play_level")

		// Extract levelNumber from URL query parameter
		levelNumberStr := c.Query("levelNumber")
		log.Println("levelNumberStr: ", levelNumberStr)

		levelNumber, err := strconv.ParseInt(levelNumberStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid levelNumber parameter"})
			return
		}

		// Define filter
		filter := bson.M{"levelNumber": levelNumber}

		var gamePlayLevel models.GamePlayLevel
		err = collection.FindOne(context.Background(), filter).Decode(&gamePlayLevel)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(404, gin.H{"error": "Game play level not found"})
				return
			}
			c.JSON(500, gin.H{"error": "Error fetching game play level"})
			return
		}

		c.JSON(200, gin.H{"gamePlayLevel": gamePlayLevel})
	}
}
