package controllers

import (

	// "log"
	// "strconv"
	// "net/http"
	// "time"
	// "fmt"

	"github.com/gofiber/fiber/v2"

	database "example/backend/database"
	"example/backend/models"
)

func GetUser(c *fiber.Ctx) error {
	userId := c.Query("userId")

	if userId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	var user models.User
	dbErr := database.MysqlDB().Where("userId = ?", userId).First(&user).Error

	if dbErr != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found!",
		})
	}

	c.Status(200).JSON(fiber.Map{"user": user})

	return nil
}
