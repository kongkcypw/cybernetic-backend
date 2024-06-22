package main

import (
	database "example/backend/database"
	routes "example/backend/routes"
	"os"

	"log"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// MessageObject Basic chat message object
type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

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

	// Initialize cors config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://gofiber.io, https://gofiber.net, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	// Define the routes
	routes.UserRoutes(app)
	routes.AuthRoutes(app)
	routes.SocketRoute(app)

	// Run the server
	app.Listen(":" + port)
}
