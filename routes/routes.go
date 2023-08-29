// routes/routes.go
package routes

import (
	"ezcom/handlers"
	"ezcom/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup initializes the routes and handlers
func Setup(router *gin.Engine) {
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Text": "Hello"})
	})
	productGroup := router.Group("/products")
	{
		productGroup.POST("", handlers.UploadImage)
		productGroup.GET("/:id", handlers.GetProductByID)
		productGroup.GET("/", handlers.GetProducts)
		productGroup.PUT("/:id", handlers.UpdateProduct)
		productGroup.DELETE("/:id", handlers.DeleteProduct)
	}
	userGroup := router.Group("/users")
	{
		userGroup.POST("", handlers.CreateMember)
		userGroup.GET("/:id", handlers.GetUser)
	}
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handlers.Singup)
		authGroup.POST("/login", handlers.Login)
		authGroup.GET("/validate", middleware.RequireAuth, handlers.Validate)
	}
}
