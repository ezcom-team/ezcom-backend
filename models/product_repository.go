package models

import (
	"context"
	"ezcom/db"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func UpdateProductQuantity(pid string) error {
	var collection = db.GetSellOrder_Collection()

	// นับ จำนวน sellOrder
	count, err := collection.CountDocuments(context.Background(), bson.M{"product_id": pid})
	if err != nil {
		return err
	}
	// หา price ที่ถูกที่สุดใน db haha
	filter := bson.M{"product_id": pid}
	options := options.FindOne().SetSort(bson.D{{Key: "price", Value: 1}})

	// Find the document with the smallest value
	var result SellOrder
	var price float64
	err = collection.FindOne(context.Background(), filter, options).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No documents found in sellOrder or sellOrder is empty")
			price = 0
		} else {
			return err
		}

	} else {
		price = result.Price
	}

	// เปลี่ยนค่า ใน Product
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{
		"$set": bson.M{
			"quantity": count,
			"price":    price,
		},
	}
	collection = db.GetProcuct_Collection()
	_, err = collection.UpdateOne(ctx, bson.M{"product_id": pid}, update)
	if err != nil {
		return err
	}
	return nil
}
