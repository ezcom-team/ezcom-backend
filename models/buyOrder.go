package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuyOrder struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Buyer_id   string             `bson:"buyer_id"`
	Buyer_name string             `bson:"buyer_name"`
	Price      float64            `bson:"price" binding:"required"`
	Product_id string             `bson:"product_id" binding:"required"`
	Condition  []string           `bson:"condition" binding:"required"`
	Color      []string           `bson:"color" binding:"required"`
	CreatedAt  time.Time          `bson:"createAt"`
}
