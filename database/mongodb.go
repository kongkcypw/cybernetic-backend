package database

import (
	"context"
	"fmt"
	"log"

	// "os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() *mongo.Client {

	// connectURL := os.Getenv("MONGODB_URL")
	connectURL := "mongodb+srv://kongkcypw:8dpPoirdcLtDk0pX@cybernatic-cluster.2f0le4p.mongodb.net/?retryWrites=true&w=majority&appName=cybernatic-cluster"
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

// var client *mongo.Client = InitMongoDB()
var Client *mongo.Client = InitMongoDB()

func MongoDB() *mongo.Client {
	return Client
}

// func OpenCollection(client *mongo.Client, databaseName, collectionName string) *mongo.Collection {
// 	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
// 	return collection
// }
