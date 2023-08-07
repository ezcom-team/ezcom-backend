package models

type Stock struct {
	Product_ID string  `bson:"product_id"`
	Price      float64 `bson:"price"`
	Owner_ID   string  `bson:"owner_id"`
}
