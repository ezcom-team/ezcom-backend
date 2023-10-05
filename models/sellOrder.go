package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SellOrder struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Seller_id  string             `bson:"seller_id"`
	Price      float64            `bson:"price"`
	Product_id string             `bson:"product_id"`
	CreatedAt  time.Time          `bson:"createAt"`
}
