package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `form:"name" binding:"required"`
	Email    string             `form:"email" binding:"required"`
	Password string             `form:"password" binding:"required"`
	Role     string             `form:"role" binding:"required"`
	Point    float64            `bson:"point"`
	File     string             `form:"file"`
}
