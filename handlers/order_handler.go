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
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateSellOrder(c *gin.Context) {
	//รับข้อมูลจาก body
	var sellOrder models.SellOrder
	if err := c.ShouldBindJSON(&sellOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "what the fuck"})
		return
	}
	fmt.Println("req")
	fmt.Println(sellOrder)

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
		"price":      bson.M{"$gte": sellOrder.Price},
		"color":      sellOrder.Color,
		"product_id": sellOrder.Product_id,
		"condition":  sellOrder.Condition,
	}

	// กำหนด options เพื่อเรียงลำดับตาม create_at ในลำดับจากน้อยไปมาก
	// options := options.FindOne().SetSort(bson.D{{Key: "createAt", Value: 1}})

	err := collection.FindOne(context.Background(), filter).Decode(&buyOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			match = false
		} else {
			panic(err)
		}

	}
	// ดึงข้อมูล product มาเพื่อ เอา P.image และ P.name
	productFound, err := models.GetProductByIdD(sellOrder.Product_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if match {
		sellOrder.Seller_id = userObj.ID.Hex()
		sellOrder.CreatedAt = time.Now()
		// create matchedOder
		var matchedOrder models.MatchedOrder
		matchedOrder.Product_img = productFound.Image
		matchedOrder.Product_name = productFound.Name
		matchedOrder.Buyer_id = buyOrder.Buyer_id
		matchedOrder.Seller_id = sellOrder.Seller_id
		matchedOrder.Color = sellOrder.Color
		matchedOrder.Condition = sellOrder.Condition
		matchedOrder.Price = sellOrder.Price
		matchedOrder.Product_id = sellOrder.Product_id
		matchedOrder.Status = "prepare"
		matchedOrder.CreatedAt = time.Now()
		collection = db.GetMatchOrder_Collection()
		result, err := collection.InsertOne(context.Background(), matchedOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			c.JSON(http.StatusOK, matchedOrder)
		}
		c.JSON(http.StatusCreated, gin.H{
			"result": result.InsertedID,
			"type":   "matchedOrder",
		})

		//delete buyorder
		var collection = db.GetBuyOrder_Collection()
		_, err = collection.DeleteOne(context.Background(), bson.M{"_id": buyOrder.ID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete buyOrder"})
			return
		}
	} else {
		sellOrder.Seller_id = userObj.ID.Hex()
		sellOrder.CreatedAt = time.Now()
		sellOrder.Seller_name = userObj.Name
		sellOrder.Product_img = productFound.Image
		sellOrder.Product_name = productFound.Name
		//สร้างข้อมูลใน DB
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		collection = db.GetSellOrder_Collection()
		result, err := collection.InsertOne(ctx, sellOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}
		// อัพเดทค่าใน product
		// productObjID, err := primitive.ObjectIDFromHex(sellOrder.Product_id)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// 	return
		// }
		err = models.UpdateProductQuantity(sellOrder.Product_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		// ส่งค่าให้ client
		c.JSON(http.StatusCreated, gin.H{
			"result": result.InsertedID,
			"type":   "sellOrder",
		})

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
		"color":      bson.M{"$in": buyOrder.Color},
		"condition":  bson.M{"$in": buyOrder.Condition},
		"price":      bson.M{"$lte": buyOrder.Price},
		"product_id": buyOrder.Product_id,
	}

	// กำหนด options เพื่อเรียงลำดับตาม create_at ในลำดับจากน้อยไปมาก
	// options := options.FindOne().SetSort(bson.D{{Key: "createAt", Value: 1}})

	err := collection.FindOne(context.Background(), filter).Decode(&sellOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			match = false
		} else {
			panic(err)
		}
	}
	fmt.Println("match is ", match)
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
	productFound, err := models.GetProductByIdD(buyOrder.Product_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println("product found is =")
	fmt.Println(productFound)
	if match {
		// create matchedOder
		var matchedOrder models.MatchedOrder
		matchedOrder.Product_img = productFound.Image
		matchedOrder.Product_name = productFound.Name
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

		//delete sellOrder
		var collection = db.GetSellOrder_Collection()
		_, err = collection.DeleteOne(context.Background(), bson.M{"_id": sellOrder.ID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sellOrder"})
			return
		}
		//update product
		// productObjID, err := primitive.ObjectIDFromHex(sellOrder.Product_id)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err})
		// 	return
		// }
		err = models.UpdateProductQuantity(buyOrder.Product_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
	} else {
		// var buyOrder models.BuyOrder
		buyOrder.Product_img = productFound.Image
		buyOrder.Product_name = productFound.Name
		buyOrder.Buyer_id = userObj.ID.Hex()
		buyOrder.Buyer_name = userObj.Name
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
func GetBuyOrdersByPID(c *gin.Context) {
	//ดึงค่า user จากใน context
	productID := c.Param("pid")

	filter := bson.M{"product_id": productID}

	// MongoDB query to find users by uid and type
	collection := db.GetBuyOrder_Collection()
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
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
func GetSellOrdersByPID(c *gin.Context) {
	//ดึงค่า user จากใน context
	productID := c.Param("pid")

	filter := bson.M{"product_id": productID}
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
