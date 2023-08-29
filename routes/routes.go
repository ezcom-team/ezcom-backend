// routes/routes.go
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup initializes the routes and handlers
func Setup(router *gin.Engine) {
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Text": "Test ja"})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Text": "Hello"})
	})
}
