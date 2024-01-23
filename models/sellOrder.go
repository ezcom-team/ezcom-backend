package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SellOrder struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Seller_id   string             `bson:"seller_id"`
	Seller_name string             `bson:"seller_name"`
	Price       float64            `bson:"price" binding:"required"`
	Product_id  string             `bson:"product_id" binding:"required"`
	Condition   string             `bson:"condition" binding:"required"`
	Color       string             `bson:"color" binding:"required"`
	CreatedAt   time.Time          `bson:"createAt"`
}
