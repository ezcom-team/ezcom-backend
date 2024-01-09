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
	Condition  string             `bson:"condition"`
	Color      string             `bson:"color"`
	CreatedAt  time.Time          `bson:"createAt"`
}
