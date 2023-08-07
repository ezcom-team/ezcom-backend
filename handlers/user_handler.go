// handlers/product_handler.go
package handlers

import (
	"context"
	"ezcom/db"
	"ezcom/models"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func CreateUser(c *gin.Context) {
	var user models.User
	log.Println("body data :", c.Request.Body)
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var collection = db.GetUser_Collection()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, result.InsertedID)
}

func UploadHandler(c *gin.Context) {
	// รับข้อมูล name และ age จาก Form Data
	name := c.PostForm("name")
	age := c.PostForm("age")

	// รับไฟล์จาก Form Data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File not found",
		})
		return
	}

	// บันทึกไฟล์ในที่เก็บที่คุณต้องการ
	// เช่นในที่นี้คือในโฟลเดอร์ uploads
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.Mkdir(uploadsDir, 0755)
	}

	filename := filepath.Join(uploadsDir, file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot save file",
		})
		return
	}
	collection := db.GetUser_Collection()
	person := models.User{
		Name: name,
		Role: age,
		File: file.Filename,
	}
	// บันทึกข้อมูลใน MongoDB
	_, err = collection.InsertOne(context.Background(), person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot save data to MongoDB",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"age":  age,
		"file": file.Filename,
	})
}

func CreateMember(c *gin.Context) {
	// รับข้อมูล name และ age จาก Form Data
	name := c.PostForm("name")
	age := c.PostForm("age")

	// รับไฟล์จาก Form Data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File not found",
		})
		return
	}

	// บันทึกไฟล์ใน GridFS
	database := db.GetDB()
	bucket, err := gridfs.NewBucket(
		database,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot connect to GridFS",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot open file",
		})
		return
	}
	defer src.Close()

	// สร้างไฟล์ใหม่ใน GridFS
	uploadStream, err := bucket.OpenUploadStream(file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot create GridFS upload stream",
		})
		return
	}
	defer uploadStream.Close()

	// คัดลอกข้อมูลจากไฟล์ที่รับมาไปยังไฟล์ใหม่ใน GridFS
	_, err = io.Copy(uploadStream, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot write file to GridFS",
		})
		return
	}

	// บันทึกข้อมูล Person ใน MongoDB
	collection := db.GetUser_Collection()
	person := models.User{
		Name: name,
		Role: age,
		File: file.Filename,
	}
	_, err = collection.InsertOne(context.Background(), person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot save data to MongoDB",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"age":  age,
		"file": file.Filename,
	})
}

func GetUser(c *gin.Context) {
	// เลือกฐานข้อมูลและคอลเล็กชันที่ต้องการเก็บข้อมูล
	collection := db.GetUser_Collection()

	// รับพารามิเตอร์ id ที่ส่งมาจาก URL
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// ดึงข้อมูล Person จาก MongoDB โดยใช้ _id เป็นเงื่อนไขในการค้นหา
	var person models.User
	filter := bson.M{"_id": objID}
	err = collection.FindOne(context.Background(), filter).Decode(&person)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Person not found",
		})
		return
	}

	// ดึงไฟล์ที่เก็บใน MongoDB GridFS
	bucket, err := gridfs.NewBucket(
		db.GetDB(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot connect to GridFS",
		})
		return
	}

	file, err := bucket.OpenDownloadStreamByName(person.File)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot open GridFS download stream",
		})
		return
	}
	defer file.Close()

	// อ่านข้อมูลจากไฟล์และส่งกลับไปยัง Client
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot read file data",
		})
		return
	}

	// ส่งข้อมูล Person กลับไปยัง Client พร้อมกับข้อมูลภาพในรูปแบบ Base64
	c.JSON(http.StatusOK, gin.H{
		"name": person.Name,
		"age":  person.Role,
		"file": data,
	})
}
