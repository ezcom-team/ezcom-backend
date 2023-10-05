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
	client               *mongo.Client
	database             *mongo.Database
	product_collection   *mongo.Collection
	user_collection      *mongo.Collection
	sellOrder_collection *mongo.Collection
)

// Connect initializes the MongoDB connection
func Connect() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://ezcom-dev:1234ezcom@cluster0.xenfcls.mongodb.net/?retryWrites=true&w=majority")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
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

	if err != nil {
		log.Fatal(err)
	}

	database = client.Database("ezcom-test")
	product_collection = database.Collection("products")
	user_collection = database.Collection("user")
	sellOrder_collection = database.Collection("sellOrder")

	log.Println("Connected to MongoDB!")
}

func GetDB() *mongo.Database {
	return database
}

func GetProcuct_Collection() *mongo.Collection {
	return product_collection
}

func GetUser_Collection() *mongo.Collection {
	return user_collection
}

func GetSellOrder_Collection() *mongo.Collection {
	return sellOrder_collection
}
