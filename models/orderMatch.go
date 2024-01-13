package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchedOrder struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Buyer_id   string             `bson:"buyer_id"`
	BuyerName  string             `bson:"buyerName"`
	Seller_id  string             `bson:"seller_id"`
	SellerName string             `bson:"sellerName"`
	Price      float64            `bson:"price"`
	Product_id string             `bson:"product_id"`
	Condition  string             `bson:"condition"`
	Color      string             `bson:"color"`
	Status     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"createAt"`
}
