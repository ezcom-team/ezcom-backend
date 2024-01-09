package handlers

import (
	"context"
	"ezcom/db"
	"ezcom/models"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateSellOrder(c *gin.Context) {
	//รับข้อมูลจาก body
	var sellOrder models.SellOrder
	if err := c.ShouldBindJSON(&sellOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	todo //check ข้อมูลฃ
	// - ดึง buyorder ที่ราคามากกว่าขึ้นไป
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

func CreateBuyOrder(c *gin.Context) {
	//รับ body
	var buyOrder models.BuyOrder
	if err := c.ShouldBindJSON(&buyOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	//เช็ค orderMatch
	// - ดึงข้อมูล sellorder ที่ถูก/น้อยกว่าทั้งหมด
	collection := db.GetSellOrder_Collection()
	var sellOrder models.SellOrder
	match := true
	filter := bson.M{"price": bson.M{"%lte": buyOrder.Price}}
	err := collection.FindOne(context.Background(), filter).Decode(&sellOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			match = false
		}
		panic(err)
	}
	//เพิ่ม ordermath or buyorder in database
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
	
	if match {
		sellOrder.Seller_id = userObj.ID.Hex()
		sellOrder.CreatedAt = time.Now()
		// create matchedOder
		var matchedOrder models.MatchedOrder
		matchedOrder.Buyer_id = userObj.ID.Hex()
		matchedOrder.Seller_id = sellOrder.Seller_id
		matchedOrder.Color = sellOrder.Color
		matchedOrder.Condition = sellOrder.Condition
		matchedOrder.Price = sellOrder.Price
		matchedOrder.Product_id = sellOrder.Product_id
		matchedOrder.CreatedAt = time.Now()
		collection = db.GetMatchOrder_Collection()
		_,err = collection.InsertOne(context.Background(),matchedOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			return
		} else {
			c.JSON(http.StatusOK,matchedOrder)
			return 
		}
	} else {
		todo// - ดึงข้อมูล user
		var buyOrder models.BuyOrder
		collection = db.GetBuyOrder_Collection()
		_,err = collection.InsertOne(context.Background(),buyOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			return
		} else {
			c.JSON(http.StatusOK,buyOrder)
			return 
		}
	}
	
	//return to client
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
