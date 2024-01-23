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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	//check ข้อมูลฃ
	// - ดึง buyorder ที่ราคามากกว่าขึ้นไป (1 of 1)
	collection := db.GetBuyOrder_Collection()
	var buyOrder models.BuyOrder
	match := true
	filter := bson.M{
		"price":      bson.M{"%gte": sellOrder.Price},
		"color":      sellOrder.Color,
		"product_id": sellOrder.Product_id,
		"condition":  sellOrder.Condition,
	}

	// กำหนด options เพื่อเรียงลำดับตาม create_at ในลำดับจากน้อยไปมาก
	options := options.FindOne().SetSort(bson.D{{Key: "createAt", Value: 1}})

	err := collection.FindOne(context.Background(), filter, options).Decode(&buyOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			match = false
		}
		panic(err)
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
		result, err := collection.InsertOne(context.Background(), matchedOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			c.JSON(http.StatusOK, matchedOrder)
		}
		c.JSON(http.StatusCreated, result)
	} else {
		var sellOrder models.SellOrder
		sellOrder.Seller_id = userObj.ID.Hex()
		sellOrder.CreatedAt = time.Now()
		//สร้างข้อมูลใน DB
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		collection = db.GetSellOrder_Collection()
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
	filter := bson.M{
		"price":      bson.M{"%lte": buyOrder.Price},
		"color":      buyOrder.Color,
		"product_id": buyOrder.Product_id,
		"condition":  buyOrder.Condition,
	}

	// กำหนด options เพื่อเรียงลำดับตาม create_at ในลำดับจากน้อยไปมาก
	options := options.FindOne().SetSort(bson.D{{Key: "createAt", Value: 1}})

	err := collection.FindOne(context.Background(), filter, options).Decode(&sellOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			match = false
		} else {
			panic(err)
		}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error can't find"})
		return
	}

	if match {
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
		result, err := collection.InsertOne(context.Background(), matchedOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		c.JSON(http.StatusCreated, gin.H{
			"result": result.InsertedID,
			"type":   "matchedOrder",
		})
	} else {
		// var buyOrder models.BuyOrder
		buyOrder.Buyer_id = userObj.ID.Hex()
		buyOrder.CreatedAt = time.Now()
		collection = db.GetBuyOrder_Collection()
		result, err := collection.InsertOne(context.Background(), buyOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		c.JSON(http.StatusCreated, gin.H{
			"result": result.InsertedID,
			"type":   "buyOrder",
		})
	}
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error 2222"})
		return
	}
	fmt.Print(userObj.ID.Hex())
	filter := bson.M{"seller_id": userObj.ID.Hex()}

	// MongoDB query to find users by uid and type
	collection := db.GetSellOrder_Collection()
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error 3333"})
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the results and decode them into the users slice
	var sellOrders []models.SellOrder
	for cursor.Next(context.TODO()) {
		var sellOrder models.SellOrder
		if err := cursor.Decode(&sellOrder); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error 44444"})
			return
		}
		sellOrders = append(sellOrders, sellOrder)
	}

	// Return filtered users as JSON
	c.JSON(http.StatusOK, sellOrders)
}
func GetBuyOrdersByUID(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error 5555"})
		return
	}
	fmt.Print(userObj.ID.Hex())
	filter := bson.M{"buyer_id": userObj.ID.Hex()}

	// MongoDB query to find users by uid and type
	collection := db.GetBuyOrder_Collection()
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error 66666"})
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the results and decode them into the users slice
	var buyOrders []models.BuyOrder
	for cursor.Next(context.TODO()) {
		var buyOrder models.BuyOrder
		if err := cursor.Decode(&buyOrder); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error 7777"})
			return
		}
		buyOrders = append(buyOrders, buyOrder)
	}

	// Return filtered users as JSON
	c.JSON(http.StatusOK, buyOrders)
}
func GetSellOrders(c *gin.Context) {

	// MongoDB query to find users by uid and type
	collection := db.GetSellOrder_Collection()
	cursor, err := collection.Find(context.TODO(), bson.M{})
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
func GetBuyOrders(c *gin.Context) {
	// MongoDB query to find users by uid and type
	collection := db.GetBuyOrder_Collection()
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the results and decode them into the users slice
	var buyOrders []models.BuyOrder
	for cursor.Next(context.TODO()) {
		var buyOrder models.BuyOrder
		if err := cursor.Decode(&buyOrder); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		buyOrders = append(buyOrders, buyOrder)
	}

	// Return filtered users as JSON
	c.JSON(http.StatusOK, buyOrders)
}
