package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgressDetail struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	MatchedOrder_id string             `bson:"matchedOrder_id"`
	Desc            string             `bson:"desc"`
	Time_duration   string             `bson:"time_duration"`
	Status          string             `bson:"status"`
	Tag_shiping     string             `bson:"tag_shiping"`
	CreatedAt       time.Time          `bson:"createAt"`
}
