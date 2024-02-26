package models

import (
	"context"
	"ezcom/db"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByIdD(uid string) (User, error) {
	userObjID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return User{}, err
	}
	collection := db.GetUser_Collection()
	var userFound User
	err = collection.FindOne(context.Background(), bson.M{"_id": userObjID}).Decode(&userFound)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No documents found in user is empty")
		} else {
			return User{}, err
		}
	}
	return userFound, nil
}
