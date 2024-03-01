package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchedOrder struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Buyer_id        string             `bson:"buyer_id"`
	BuyerName       string             `bson:"buyerName"`
	Buyer_img       string             `bson:"buyer_img"`
	Seller_id       string             `bson:"seller_id"`
	SellerName      string             `bson:"sellerName"`
	Seller_img      string             `bson:"seller_img"`
	Price           float64            `bson:"price"`
	Product_id      string             `bson:"product_id"`
	Product_name    string             `bson:"product_name"`
	Product_img     string             `bson:"product_img"`
	Condition       string             `bson:"condition"`
	Color           string             `bson:"color"`
	Status          string             `bson:"status"`
	CreatedAt       time.Time          `bson:"createdAt"`
	Verify          string             `bson:"verify"`
	Received        string             `bson:"received"`
	DesAdd          string             `bson:"desAdd"`
	DesPhone        string             `bson:"desPhone"`
	Tracking_Number string             `bson:"tracking_Number"`
	PaymentStatus   string             `bson:"paymentStatus"`
}
