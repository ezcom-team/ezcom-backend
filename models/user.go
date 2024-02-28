package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" binding:"required"`
	Email     string             `bson:"email" binding:"required"`
	Password  string             `bson:"password" binding:"required"`
	Role      string             `bson:"role" binding:"required"`
	Point     float64            `bson:"point"`
	File      string             `bson:"file"`
	CreatedAt time.Time          `bson:"createdAt"`
}
