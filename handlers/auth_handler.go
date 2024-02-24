package handlers

import (
	"context"
	"ezcom/db"
	"ezcom/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type RequestBody struct {
	Email    string
	Password string
}

func Singup(c *gin.Context) {
	// initial ctx
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	// bind request.body with user
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// user.Name = c.PostForm("name")
	// user.Email = c.PostForm("email")
	// user.Password = c.PostForm("password")
	// user.Role = c.PostForm("role")
	// fmt.Print("user data => ")
	fmt.Print(user.Name, user.Email, user.Password, user.Role)

	// if err := c.BindJSON(&user); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
	// 	return
	// }
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hash)

	// store file
	haveFile := true
	file, err := c.FormFile("file")
	if err != nil {
		haveFile = false
	}
	// file = user.File
	if haveFile {
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

		user.File = "https://firebasestorage.googleapis.com/v0/b/ezcom-eaa21.appspot.com/o/" + imagePath + "?alt=media"
	}
	// create the user
	collection := db.GetUser_Collection()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't create user"})
		return
	}
	// response
	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	// init ctx
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	// get user from body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// look up form db
	collection := db.GetUser_Collection()
	filter := bson.M{"email": body.Email}
	log.Print(body)
	var found models.User
	err := collection.FindOne(ctx, filter).Decode(&found)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password 1"})
		return
	}
	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password 2"})
		return
	}
	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": found.ID,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}
	// set token in cookie
	// c.SetSameSite(http.SameSiteNoneMode)
	// c.SetCookie("Authorization", tokenString, 120000, "/", "", true, false)

	// sent tokenString
	c.JSON(http.StatusOK, gin.H{
		"user":   found,
		"token":  tokenString,
		"userID": found.ID,
	})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I'm login"})
}
