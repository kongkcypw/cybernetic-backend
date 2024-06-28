package main

import (
	database "example/backend/database"
	routes "example/backend/routes"
	"net/http"
	"os"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
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

	// Create a new Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(GinMiddleware("http://localhost:5173"))

	// Socket.io server
	socketServer := routes.SocketServerRoute()
	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer socketServer.Close()
	router.GET("/socket.io/*any", gin.WrapH(socketServer))
	router.POST("/socket.io/*any", gin.WrapH(socketServer))

	// Register the routes
	routes.UserRoutes(router)
	routes.AuthRoutes(router)

	if err := router.Run(":" + port); // router.Run("localhost:" + port)
	err != nil {
		log.Fatal("failed run app: ", err)
	}
}
