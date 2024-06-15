package main

import (
	database "example/backend/database"
	routes "example/backend/routes"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// Create a new Gin router
	router := gin.New()
	router.Use(gin.Logger())

	// Define the routes
	routes.UserRoutes(router)
	routes.AuthRoutes(router)

	// Run the server
	router.Run("localhost:" + port)
}
