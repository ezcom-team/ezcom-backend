// models/product.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product struct represents a product in the database
type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Price    float64            `bson:"price"`
	Quantity int                `bson:"quantity"`
}
