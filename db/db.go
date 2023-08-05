// db/db.go
package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
)

// Connect initializes the MongoDB connection
func Connect() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://ezcom-dev:1234ezcom@cluster0.xenfcls.mongodb.net/?retryWrites=true&w=majority")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	database = client.Database("ezcom-test")
	collection = database.Collection("products")

	log.Println("Connected to MongoDB!")
}

func GetCollection() *mongo.Collection {
	return collection
}
