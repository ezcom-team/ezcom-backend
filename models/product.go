// models/product.go
package models

import (
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product struct represents a product in the database
type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" binding:"required"`
	Desc      string             `bson:"desc"`
	Price     float64            `bson:"price" binding:"required"`
	Image     multipart.File     `bson:"image"`
	ImagePath string             `bson:"imagePath" `
	Quantity  int64              `bson:"quantity"`
	Type      string             `bson:"type" binding:"required"`
	Color     []string           `bson:"color" binding:"required"`
	Specs     string             `bson:"specs"`
}
