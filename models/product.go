// models/product.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product struct represents a product in the database
type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" binding:"required"`
	Desc      string             `bson:"desc"`
	Price     float64            `bson:"price"`
	Image     string             `bson:"image"`
	Quantity  int64              `bson:"quantity"`
	Type      string             `bson:"type"`
	Color     []string           `bson:"color"`
	Specs     string             `bson:"specs"`
	CreatedAt time.Time          `bson:"createAt"`
}
