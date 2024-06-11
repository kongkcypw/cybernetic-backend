package main

import (
	orm "example/backend/config"
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
	orm.InitDB()

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	routes.AuthRoutes(router)

	router.Run("localhost:" + port)
}
