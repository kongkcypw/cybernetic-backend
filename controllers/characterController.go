package controllers

import (
	// "fmt"
	// "os"
	// "strconv"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	database "example/backend/database"
	"example/backend/models"
)

func GetCharacter(c *fiber.Ctx) error {
	userId := c.Query("userId")

	if userId == "" {
		c.Status(400).JSON(fiber.Map{
			"error": "userId is required",
		})
		return fiber.NewError(400, "userId is required")
	}

	client := database.MongoDB()
	if client == nil {
		c.Status(500).JSON(fiber.Map{
			"error": "Database connection is not initialized",
		})
		return fiber.NewError(500, "Database connection is not initialized")
	}

	collection := database.MongoDBOpenCollection(client, "cybernetic", "user_character")

	var character models.Character
	err := collection.FindOne(context.Background(), bson.M{"userId": userId}).Decode(&character)
	if err != nil {
		c.Status(404).JSON(fiber.Map{
			"error": "Character not found!",
		})
		return fiber.NewError(404, "Character not found!")
	}

	c.Status(200).JSON(fiber.Map{"character": character})
	return nil
}
