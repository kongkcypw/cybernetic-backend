package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var Client *mongo.Client

func InitMongoDB() *mongo.Client {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectURL := os.Getenv("MONGODB_URL")
	if connectURL == "" {
		log.Fatal("MONGODB_URL environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(connectURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		return nil
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Printf("Failed to ping MongoDB: %v\n", err)
		return nil
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

var Client *mongo.Client = InitMongoDB()

func MongoDB() *mongo.Client {
	return Client
}

func MongoDBOpenCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
	return collection
}
