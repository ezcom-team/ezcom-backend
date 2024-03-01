// handlers/product_handler.go
package handlers

import (
	"context"
	"ezcom/db"
	"ezcom/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

func CreateProduct(c *gin.Context) {
	var product models.Product
	product.Name = c.PostForm("name")
	product.Type = c.PostForm("type")
	product.Desc = c.PostForm("desc")
	color, ok := c.GetPostFormArray("color")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bug color"})
		return
	}
	product.Color = color
	// var product models.Product
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Print("formdata is :")
	fmt.Print(product)
	fmt.Print(product.Name, product.Type)
	priceStr := c.PostForm("price")
	priceFloat, _ := strconv.ParseFloat(priceStr, 64)
	product.Price = priceFloat
	quantityStr := c.PostForm("quantity")
	quantityInt, _ := strconv.ParseInt(quantityStr, 10, 64)
	product.Quantity = int64(quantityInt)
	product.ID = primitive.NewObjectID()
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
	product.CreatedAt = time.Now()

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
		specs.Width = c.PostForm("width") // store product in database
		specs.Height = c.PostForm("height")
		specs.DPI = c.PostForm("dpi")
		specs.PID = product.ID.Hex()
		specs.Type = product.Type
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
		specs.PID = product.ID.Hex()
		specs.Type = product.Type
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
	} else if product.Type == "headset" {
		var specs models.HeadsetSpecs
		specs.Headset_Type = c.PostForm("headset_type")
		specs.Cable_Length = c.PostForm("cable_length")
		specs.Connection = c.PostForm("connection")
		specs.Microphone = c.PostForm("microphone")
		specs.Noise_Cancelling = c.PostForm("noise_cancelling")
		specs.Weight = c.PostForm("weight")
		specs.PID = product.ID.Hex()
		specs.Type = product.Type
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
	} else if product.Type == "mousePad" {
		var specs models.MousePad
		specs.Height = c.PostForm("height")
		specs.Thickness = c.PostForm("thickness")
		specs.Material = c.PostForm("material")
		specs.Length = c.PostForm("length")
		specs.Stitched_edges = c.PostForm("stitched_edges")
		specs.Glide = c.PostForm("glide")
		specs.PID = product.ID.Hex()
		specs.Type = product.Type
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
	} else if product.Type == "GPU" {
		var specs models.GPU
		specs.NVIDIA_CUDA_Cores = c.PostForm("nvidia_cuda_cores")
		specs.Memory_Size = c.PostForm("memory_size")
		specs.Boost_Clock = c.PostForm("boost_clock")
		specs.Memory_Type = c.PostForm("memory_type")
		specs.PID = product.ID.Hex()
		specs.Type = product.Type
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
	} else if product.Type == "CPU" {
		var specs models.CPU
		specs.Socket = c.PostForm("socket")
		specs.Threads = c.PostForm("threads")
		specs.Core_Speed_Base = c.PostForm("core_speed_base")
		specs.Cores = c.PostForm("cores")
		specs.TDP = c.PostForm("TDP")
		specs.Core_Speed_Boost = c.PostForm("core_speed_boost")
		specs.PID = product.ID.Hex()
		specs.Type = product.Type
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
	} else if product.Type == "monitor" {
		var specs models.Monitor
		specs.Size = c.PostForm("size")
		specs.Aspect_Ratio = c.PostForm("aspect_ratio")
		specs.G_Sync = c.PostForm("g_sync")
		specs.Panel_Tech = c.PostForm("panel_tech")
		specs.Resolution = c.PostForm("resolution")
		specs.Refresh_Rate = c.PostForm("refresh_rate")
		specs.FreeSync = c.PostForm("free_sync")
		specs.PID = product.ID.Hex()
		specs.Type = product.Type
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

}

func GetProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the MongoDB collection
	collection := db.GetProcuct_Collection()

	// Find all products in the collection

	// "Failed to retrieve products"
	// cursor, err := collection.Find(ctx, bson.M{})
	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))

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
	// var spec interface{} // ประกาศตัวแปร spec ไว้นอก switch

	switch specType {
	case "mouse":
		var spec models.MouseSpecs
		fmt.Println("mouse")
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
		fmt.Println(spec)

		c.JSON(http.StatusOK, spec)
	case "keyboard":
		var spec models.KeyBoardSpecs
		fmt.Println("keyBoard")
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
		fmt.Println(spec)

		c.JSON(http.StatusOK, spec)
	case "headset":
		var spec models.HeadsetSpecs
		fmt.Println("headset")
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
		fmt.Println(spec)

		c.JSON(http.StatusOK, spec)
	case "mousePad":
		var spec models.MousePad
		fmt.Println("mousePad")
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
		fmt.Println(spec)

		c.JSON(http.StatusOK, spec)
	case "GPU":
		var spec models.GPU
		fmt.Println("GPU")
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
		fmt.Println(spec)

		c.JSON(http.StatusOK, spec)
	case "CPU":
		var spec models.CPU
		fmt.Println("CPU")
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
		fmt.Println(spec)

		c.JSON(http.StatusOK, spec)
	case "monitor":
		var spec models.Monitor
		fmt.Println("monitor")
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
		fmt.Println(spec)

		c.JSON(http.StatusOK, spec)
	default:
		c.JSON(http.StatusBadRequest, "don't have spec type yet")
	}
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
	product.Name = c.PostForm("name")
	product.Desc = c.PostForm("desc")
	product.Type = c.PostForm("type")
	product.Color = c.PostFormArray("color")
	product.Specs = c.PostForm("specs")
	// store file
	updataImage := true
	file, err := c.FormFile("image")
	if err != nil {
		updataImage = false
	}
	// if file != nil {
	// 	updataImage = false
	// }
	if updataImage {

		foundProduct, err := models.GetProductByIdD(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		// delete foundProduct.Image
		client, err := storage.NewClient(ctx, option.WithCredentialsFile("ezcom-eaa21-firebase-adminsdk-9zpt0-d8e4765278.json"))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		// ชื่อของ bucket ที่เก็บไฟล์
		bucketName := "ezcom-eaa21.appspot.com"

		// ชื่อของไฟล์ที่ต้องการลบ

		// Example string
		path := foundProduct.Image

		// Split the string using "/"
		parts := strings.Split(path, "/")

		// Print the last element
		lastIndex := len(parts) - 1
		parts1 := strings.Split(parts[lastIndex], "?")
		fmt.Println(parts1[0])
		fileName := parts1[0]

		// ลบไฟล์
		err = client.Bucket(bucketName).Object(fileName).Delete(ctx)
		if err != nil {
			log.Fatalf("Failed to delete object: %v", err)
		}

		fmt.Println("Object deleted successfully")

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
		fmt.Print("product image : ", product.Image)
		product.CreatedAt = foundProduct.CreatedAt
	} else {
		foundProduct, err := models.GetProductByIdD(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		product.Image = foundProduct.Image
		product.CreatedAt = foundProduct.CreatedAt
	}

	if err := c.Bind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shouldbind error"})
		return
	}

	update := bson.M{
		"$set": product, // ใช้ struct ที่ได้รับเป็นค่าในการอัปเดตทุกฟิลด์
	}
	var collection = db.GetProcuct_Collection()
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	if product.Type == "mouse" {
		var specs models.MouseSpecs
		specs.Type = c.PostForm("type")
		specs.PID = c.PostForm("pID")
		specs.Sensor = c.PostForm("sensor")
		specs.ButtonSwitch = c.PostForm("buttonSwitch")
		specs.Connection = c.PostForm("connection")
		specs.Length = c.PostForm("length")
		specs.Weight = c.PostForm("weight")
		specs.PollingRate = c.PostForm("pollingRate")
		specs.ButtonForce = c.PostForm("buttonForce")
		specs.Shape = c.PostForm("shape")
		specs.Height = c.PostForm("height")
		specs.Width = c.PostForm("width") // store product in database
		specs.DPI = c.PostForm("dpi")     // store product in database
		update := bson.M{
			"$set": specs, // ใช้ struct ที่ได้รับเป็นค่าในการอัปเดตทุกฟิลด์
		}
		collection = db.GetSpecs_Collection()
		specID, err := primitive.ObjectIDFromHex(product.Specs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
			return
		}
		result, err := collection.UpdateOne(ctx, bson.M{"_id": specID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specs"})
			return
		}
		c.JSON(http.StatusCreated, result.UpsertedID)
	} else if product.Type == "keyboard" {
		var specs models.KeyBoardSpecs
		specs.Type = c.PostForm("type")
		specs.PID = c.PostForm("pID")
		specs.Form_Factor = c.PostForm("form_factor")
		specs.PCB = c.PostForm("PCB")
		specs.RGB = c.PostForm("RGB")
		specs.Switches = c.PostForm("switches")
		specs.Length = c.PostForm("length")
		specs.Weight = c.PostForm("weight")
		specs.Height = c.PostForm("height")
		specs.Width = c.PostForm("width")

		update := bson.M{
			"$set": specs, // ใช้ struct ที่ได้รับเป็นค่าในการอัปเดตทุกฟิลด์
		}
		collection = db.GetSpecs_Collection()
		specID, err := primitive.ObjectIDFromHex(product.Specs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
			return
		}
		result, err := collection.UpdateOne(ctx, bson.M{"_id": specID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specs"})
			return
		}
		c.JSON(http.StatusCreated, result.UpsertedID)
	} else if product.Type == "headset" {
		var specs models.HeadsetSpecs
		specs.Type = c.PostForm("type")
		specs.PID = c.PostForm("pID")
		specs.Headset_Type = c.PostForm("headset_type")
		specs.Cable_Length = c.PostForm("cable_length")
		specs.Connection = c.PostForm("connection")
		specs.Microphone = c.PostForm("microphone")
		specs.Noise_Cancelling = c.PostForm("noise_cancelling")
		specs.Weight = c.PostForm("weight")

		update := bson.M{
			"$set": specs, // ใช้ struct ที่ได้รับเป็นค่าในการอัปเดตทุกฟิลด์
		}
		collection = db.GetSpecs_Collection()
		specID, err := primitive.ObjectIDFromHex(product.Specs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
			return
		}
		result, err := collection.UpdateOne(ctx, bson.M{"_id": specID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specs"})
			return
		}
		c.JSON(http.StatusCreated, result.UpsertedID)
	} else if product.Type == "mousePad" {
		var specs models.MousePad
		specs.Type = c.PostForm("type")
		specs.PID = c.PostForm("pID")
		specs.Height = c.PostForm("height")
		specs.Thickness = c.PostForm("thickness")
		specs.Material = c.PostForm("material")
		specs.Length = c.PostForm("length")
		specs.Stitched_edges = c.PostForm("stitched_edges")
		specs.Glide = c.PostForm("glide")
		update := bson.M{
			"$set": specs, // ใช้ struct ที่ได้รับเป็นค่าในการอัปเดตทุกฟิลด์
		}
		collection = db.GetSpecs_Collection()
		specID, err := primitive.ObjectIDFromHex(product.Specs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
			return
		}
		result, err := collection.UpdateOne(ctx, bson.M{"_id": specID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specs"})
			return
		}
		c.JSON(http.StatusCreated, result.UpsertedID)
	} else if product.Type == "GPU" {
		var specs models.GPU
		specs.Type = c.PostForm("type")
		specs.PID = c.PostForm("pID")
		specs.NVIDIA_CUDA_Cores = c.PostForm("nvidia_cuda_cores")
		specs.Memory_Size = c.PostForm("memory_size")
		specs.Boost_Clock = c.PostForm("boost_clock")
		specs.Memory_Type = c.PostForm("memory_type")
		update := bson.M{
			"$set": specs, // ใช้ struct ที่ได้รับเป็นค่าในการอัปเดตทุกฟิลด์
		}
		collection = db.GetSpecs_Collection()
		specID, err := primitive.ObjectIDFromHex(product.Specs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
			return
		}
		result, err := collection.UpdateOne(ctx, bson.M{"_id": specID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specs"})
			return
		}
		c.JSON(http.StatusCreated, result.UpsertedID)
	} else if product.Type == "CPU" {
		var specs models.CPU
		specs.Type = c.PostForm("type")
		specs.PID = c.PostForm("pID")
		specs.Socket = c.PostForm("socket")
		specs.Threads = c.PostForm("threads")
		specs.Core_Speed_Base = c.PostForm("core_speed_base")
		specs.Cores = c.PostForm("cores")
		specs.TDP = c.PostForm("TDP")
		specs.Core_Speed_Boost = c.PostForm("core_speed_boost")
		update := bson.M{
			"$set": specs, // ใช้ struct ที่ได้รับเป็นค่าในการอัปเดตทุกฟิลด์
		}
		collection = db.GetSpecs_Collection()
		specID, err := primitive.ObjectIDFromHex(product.Specs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
			return
		}
		result, err := collection.UpdateOne(ctx, bson.M{"_id": specID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specs"})
			return
		}
		c.JSON(http.StatusCreated, result.UpsertedID)
	} else if product.Type == "monitor" {
		var specs models.Monitor
		specs.Type = c.PostForm("type")
		specs.PID = c.PostForm("pID")
		specs.Size = c.PostForm("size")
		specs.Aspect_Ratio = c.PostForm("aspect_ratio")
		specs.G_Sync = c.PostForm("g_sync")
		specs.Panel_Tech = c.PostForm("panel_tech")
		specs.Resolution = c.PostForm("resolution")
		specs.Refresh_Rate = c.PostForm("refresh_rate")
		specs.FreeSync = c.PostForm("free_sync")
		update := bson.M{
			"$set": specs, // ใช้ struct ที่ได้รับเป็นค่าในการอัปเดตทุกฟิลด์
		}
		collection = db.GetSpecs_Collection()
		specID, err := primitive.ObjectIDFromHex(product.Specs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spec ID"})
			return
		}
		result, err := collection.UpdateOne(ctx, bson.M{"_id": specID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specs"})
			return
		}
		c.JSON(http.StatusCreated, result.UpsertedID)
	}
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
	foundProduct, err := models.GetProductByIdD(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if foundProduct.Image != "" {
		// delete foundProduct.Image
		client, err := storage.NewClient(ctx, option.WithCredentialsFile("ezcom-eaa21-firebase-adminsdk-9zpt0-d8e4765278.json"))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		// ชื่อของ bucket ที่เก็บไฟล์
		bucketName := "ezcom-eaa21.appspot.com"

		// ชื่อของไฟล์ที่ต้องการลบ

		// Example string
		path := foundProduct.Image

		// Split the string using "/"
		parts := strings.Split(path, "/")

		// Print the last element
		lastIndex := len(parts) - 1
		parts1 := strings.Split(parts[lastIndex], "?")
		fmt.Println(parts1[0])
		fileName := parts1[0]

		// ลบไฟล์
		err = client.Bucket(bucketName).Object(fileName).Delete(ctx)
		if err != nil {
			log.Fatalf("Failed to delete object: %v", err)
		}

		fmt.Println("Object deleted successfully")
	}

	var collection = db.GetSpecs_Collection()
	objSpecsID, err := primitive.ObjectIDFromHex(foundProduct.Specs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objSpecsID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete specs"})
		return
	}

	collection = db.GetProcuct_Collection()
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.Status(http.StatusOK)
}

func GetSpecs(c *gin.Context) {
	specType := c.Param("type")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch specType {
	case "mouse":
		var collection = db.GetSpecs_Collection()
		cursor, err := collection.Find(ctx, bson.M{"type": specType})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cursor.Close(ctx)
		var products []models.MouseSpecs

		for cursor.Next(ctx) {
			var product models.MouseSpecs
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
	case "keyboard":
		var collection = db.GetSpecs_Collection()
		cursor, err := collection.Find(ctx, bson.M{"type": specType})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cursor.Close(ctx)
		var products []models.KeyBoardSpecs

		for cursor.Next(ctx) {
			var product models.KeyBoardSpecs
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
	case "headset":
		var collection = db.GetSpecs_Collection()
		cursor, err := collection.Find(ctx, bson.M{"type": specType})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cursor.Close(ctx)
		var products []models.HeadsetSpecs

		for cursor.Next(ctx) {
			var product models.HeadsetSpecs
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
	case "mousePad":
		var collection = db.GetSpecs_Collection()
		cursor, err := collection.Find(ctx, bson.M{"type": specType})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cursor.Close(ctx)
		var products []models.MousePad

		for cursor.Next(ctx) {
			var product models.MousePad
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
	case "GPU":
		var collection = db.GetSpecs_Collection()
		cursor, err := collection.Find(ctx, bson.M{"type": specType})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cursor.Close(ctx)
		var products []models.GPU

		for cursor.Next(ctx) {
			var product models.GPU
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
	case "CPU":
		var collection = db.GetSpecs_Collection()
		cursor, err := collection.Find(ctx, bson.M{"type": specType})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cursor.Close(ctx)
		var products []models.CPU

		for cursor.Next(ctx) {
			var product models.CPU
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
	case "monitor":
		var collection = db.GetSpecs_Collection()
		cursor, err := collection.Find(ctx, bson.M{"type": specType})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cursor.Close(ctx)
		var products []models.Monitor

		for cursor.Next(ctx) {
			var product models.Monitor
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
	default:
		c.JSON(http.StatusBadRequest, "don't have spec type yet")
	}
}
