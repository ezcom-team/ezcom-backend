package handlers

import (
	"context"
	"ezcom/db"
	"ezcom/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSellOrder(c *gin.Context) {
	//รับข้อมูลจาก body
	var sellOrder models.SellOrder
	if err := c.ShouldBindJSON(&sellOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	//ดึงค่า user จากใน context
	user, exists := c.Get("user")
	if !exists {
		// ไม่พบค่า "user" ใน context
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	// แปลง user เป็น models.User
	userObj, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	sellOrder.Seller_id = userObj.ID.Hex()
	sellOrder.CreatedAt = time.Now()

	//สร้างข้อมูลใน DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var collection = db.GetSellOrder_Collection()
	result, err := collection.InsertOne(ctx, sellOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	// update product increase Quantity
	var foundProduct models.Product
	productObjID, err := primitive.ObjectIDFromHex(sellOrder.Product_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	foundProduct, err = models.GetProduct(productObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find product"})
		return
	}
	if sellOrder.Price < foundProduct.Price {
		err = models.UpdateProductQuantity(productObjID, foundProduct.Quantity+1, sellOrder.Price)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
	} else {
		err = models.UpdateProductQuantity(productObjID, foundProduct.Quantity+1, foundProduct.Price)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
	}

	//ส่งค่ากลับไปให้ client
	c.JSON(http.StatusCreated, result)
}

func GetSellOrdersByUID(c *gin.Context) {
	//ดึงค่า user จากใน context
	user, exists := c.Get("user")
	if !exists {
		// ไม่พบค่า "user" ใน context
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	// แปลง user เป็น models.User
	userObj, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	fmt.Print(userObj.ID.Hex())
	filter := bson.M{"seller_id": userObj.ID.Hex()}

	// MongoDB query to find users by uid and type
	collection := db.GetSellOrder_Collection()
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the results and decode them into the users slice
	var sellOrders []models.SellOrder
	for cursor.Next(context.TODO()) {
		var sellOrder models.SellOrder
		if err := cursor.Decode(&sellOrder); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		sellOrders = append(sellOrders, sellOrder)
	}

	// Return filtered users as JSON
	c.JSON(http.StatusOK, sellOrders)
}
