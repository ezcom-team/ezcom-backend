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
	router.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Text": "Hello welcome to ezcom backend"})
	})
	productGroup := router.Group("/products")
	{
		productGroup.POST("", handlers.CreateProduct)
		productGroup.GET("/:id", handlers.GetProductByID)
		productGroup.GET("/", handlers.GetProducts)
		productGroup.PUT("/:id", handlers.UpdateProduct)
		productGroup.DELETE("/:id", handlers.DeleteProduct)
	}
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handlers.Singup)
		authGroup.POST("/login", handlers.Login)
		authGroup.GET("/validate", middleware.RequireAuth, handlers.Validate)
	}
	sellorderGroup := router.Group("/sellOrder")
	{
		sellorderGroup.POST("", middleware.RequireAuth, handlers.CreateSellOrder)
		sellorderGroup.GET("", middleware.RequireAuth, handlers.GetSellOrdersByUID)
	}
	specsGroup := router.Group("/specs")
	{
		specsGroup.GET("/mouse",middleware.RequireAuth,)
		specsGroup.GET("/mousepad",middleware.RequireAuth,)
		specsGroup.GET("/cpu",middleware.RequireAuth,)
	}
}
