package models

import (
	"context"
	"ezcom/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetProduct(objID primitive.ObjectID) (Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var product Product
	var collection = db.GetProcuct_Collection()
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return Product{}, err
	}
	return product, nil

}

func UpdateProductQuantity(objID primitive.ObjectID, quantity int64, price float64) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{
		"$set": bson.M{
			"quantity": quantity,
			"price":    price,
		},
	}
	var collection = db.GetProcuct_Collection()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	return nil
}
