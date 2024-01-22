package middleware

import (
	"context"
	"ezcom/db"
	"ezcom/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequireAuth(c *gin.Context) {
	// set ctx
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	//---------
	// Get the token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	//---------
	// Get the cookie off req
	// tokenString, err := c.Cookie("Authorization")
	// tokenString, err := c.Cookie("Authorization")
	// tokenStringFromReq, _ := c.Request.Cookie("Authorization")
	// cookieValue := c.Request.Header

	// fmt.Println("tokenString")
	// fmt.Println(tokenStringFromReq)
	// fmt.Println(cookieValue)
	// fmt.Println(err)

	// if !tokenString {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "tokenString"})
	// 	return
	// }
	// Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "exp"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Find the user with token sub
		var user models.User
		collection := db.GetUser_Collection()
		objId, err := primitive.ObjectIDFromHex(claims["sub"].(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}
		err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "findOne"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Attach to req
		c.Set("user", user)
		// Continue
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "else"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}
