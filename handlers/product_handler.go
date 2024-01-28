// handlers/product_handler.go
package handlers

import (
	"context"
	"ezcom/db"
	"ezcom/models"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type createProductDTO struct {
	Name  string
	Desc  string
	Image string
	Type  string
	Color []string
	Specs string
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	product.Name = c.PostForm("name")
	product.Type = c.PostForm("type")
	product.Desc = c.PostForm("desc")
	product.Color = c.PostFormArray("color")
	priceStr := c.PostForm("price")
	priceFloat, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		// Handle the error, possibly return an error response
	}
	product.Price = priceFloat
	quantityStr := c.PostForm("quantity")
	quantityInt, err := strconv.ParseInt(quantityStr, 10, 64)
	if err != nil {
		// Handle the error, possibly return an error response
	}
	product.Quantity = int64(quantityInt)
	// store file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File not found",
		})
		return
	}
	imagePath := file.Filename

	bucket := "ezcom-eaa21.appspot.com"

	wc := db.Storage.Bucket(bucket).Object(imagePath).NewWriter(context.Background())
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer src.Close()

	_, err = io.Copy(wc, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := wc.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to close the file writer",
		})
		return
	}

	product.Image = "https://firebasestorage.googleapis.com/v0/b/ezcom-eaa21.appspot.com/o/" + imagePath + "?alt=media"
	// check and set specs

	var specs models.MouseSpecs
	if product.Type == "mouse" {
		var specs models.MouseSpecs
		specs.Sensor = c.PostForm("sensor")
		specs.ButtonSwitch = c.PostForm("buttonSwitch")
		specs.Connection = c.PostForm("connection")
		specs.Length = c.PostForm("length")
		specs.Weight = c.PostForm("weight")
		specs.PollingRate = c.PostForm("pollingRate")
		specs.ButtonForce = c.PostForm("buttonForce")
		specs.Shape = c.PostForm("shape")
		specs.Height = c.PostForm("height")
		specs.Width = c.PostForm("width")
	} else if product.Type == "keyboard" {
		var specs models.KeyBoardSpecs
		specs.Form_Factor = c.PostForm("form_factor")
		specs.PCB = c.PostForm("PCB")
		specs.RGB = c.PostForm("RGB")
		specs.Switches = c.PostForm("switches")
		specs.Length = c.PostForm("length")
		specs.Weight = c.PostForm("weight")
		specs.Height = c.PostForm("height")
		specs.Width = c.PostForm("width")
	} else if product.Type == "headset" {
		var specs models.HeadsetSpecs
		specs.Headset_Type = c.PostForm("headset_type")
		specs.Cable_Length = c.PostForm("cable_length")
		specs.Connection = c.PostForm("connection")
		specs.Microphone = c.PostForm("microphone")
		specs.Noise_Cancelling = c.PostForm("noise_cancelling")
		specs.Weight = c.PostForm("weight")
	}
	// store product in database
	var specsCollection = db.GetSpecs_Collection()
	specsResult, err := specsCollection.InsertOne(context.Background(), specs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create specs"})
		return
	}
	product.Specs = specsResult.InsertedID.(primitive.ObjectID).Hex()
	c.JSON(http.StatusCreated, specsResult.InsertedID)
	var collection = db.GetProcuct_Collection()
	result, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, result.InsertedID)
}

func GetProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the MongoDB collection
	collection := db.GetProcuct_Collection()

	// Find all products in the collection

	// "Failed to retrieve products"
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer cursor.Close(ctx)
	// Prepare a slice to hold the products
	var products []models.Product

	// Iterate through the cursor and decode each product
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode product data"})
			return
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	// Return the products as JSON response
	c.JSON(http.StatusOK, products)
}

// GetProductByID retrieves a product by its ID
func GetProductByID(c *gin.Context) {
	productID := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	var collection = db.GetProcuct_Collection()

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

type findTypeById struct {
}

// GetProductByID with specs
func GetSpecByID(c *gin.Context) {
	specID := c.Param("id")
	specType := c.Param("type")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(specID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	var spec interface{} // ประกาศตัวแปร spec ไว้นอก switch

	switch specType {
	case "Mouse":
		spec = models.MouseSpecs{}
	case "Keyboard":
		spec = models.KeyBoardSpecs{}
	case "Headset":
		spec = models.HeadsetSpecs{}
	default:
		spec = models.MouseSpecs{}
	}

	var collection = db.GetSpecs_Collection()

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&spec)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Spec not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve spec"})
		}
		return
	}

	c.JSON(http.StatusOK, spec)
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":  product.Name,
			"price": product.Price,
			"file":  product.Image,
		},
	}
	var collection = db.GetProcuct_Collection()
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteProduct deletes a product by its ID
func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var collection = db.GetProcuct_Collection()
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.Status(http.StatusNoContent)
}
