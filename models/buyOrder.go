package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuyOrder struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Buyer_id     string             `bson:"buyer_id"`
	Buyer_name   string             `bson:"buyer_name"`
	Buyer_img   string             `bson:"buyer_img"`
	Price        float64            `bson:"price" binding:"required"`
	Product_id   string             `bson:"product_id" binding:"required"`
	Product_name string             `bson:"product_name"`
	Product_img  string             `bson:"product_img"`
	Condition    []string           `bson:"condition" binding:"required"`
	Color        []string           `bson:"color" binding:"required"`
	CreatedAt    time.Time          `bson:"createAt"`
	Verify       string             `bson:"verify"`
}
