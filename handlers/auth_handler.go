package handlers

import (
	"context"
	"ezcom/db"
	"ezcom/models"
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
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}
	// create the user
	user.Password = string(hash)
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
	// // set in token
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	// return response
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I'm login"})
}
