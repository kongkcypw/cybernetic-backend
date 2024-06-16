package main

import (
	database "example/backend/database"
	routes "example/backend/routes"
	"os"

	"log"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load the .env file
	err := godotenv.Load()

	port := os.Getenv("PORT")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	database.InitMySQL()
	database.InitMongoDB()

	// Create a new Fiber app
	app := fiber.New()

	// Use the logger middleware
	app.Use(logger.New())

	// Define the routes
	routes.UserRoutes(app)
	routes.AuthRoutes(app)

	// Run the server
	app.Listen(":" + port)
}
