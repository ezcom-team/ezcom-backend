// models/product.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product struct represents a product in the database
type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `form:"name" binding:"required"`
	Desc     string             `form:"desc" `
	Price    float64            `form:"price" binding:"required"`
	Image    string             `form:"image" `
	Quantity int64              `form:"quantity"`
	Type     string             `form:"type" binding:"required"`
	Color    []string           `form:"color" binding:"required"`
	Specs    string             `form:"specs"`
}
