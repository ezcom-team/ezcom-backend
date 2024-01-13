package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuyOrder struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Buyer_id  string             `bson:"buyer_id"`
	Price      float64            `bson:"price"`
	Product_id string             `bson:"product_id"`
	Condition  []string             `bson:"condition"`
	Color      []string             `bson:"color"`
	CreatedAt  time.Time          `bson:"createAt"`
}
