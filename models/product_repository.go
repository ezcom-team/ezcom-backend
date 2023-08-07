package models

import (
	"context"
	"ezcom/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the MongoDB collection
	collection := db.GetProcuct_Collection()

	// Find all products in the collection
	// "Failed to retrieve products"
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	// Prepare a slice to hold the products
	var products []Product

	// Iterate through the cursor and decode each product
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
