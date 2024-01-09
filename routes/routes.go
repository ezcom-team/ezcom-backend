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
		productGroup.GET("/", handlers.GetProducts)
		productGroup.GET("/:id", handlers.GetProductByID)
		productGroup.GET("/spec/:id", handlers.GetSpecByID)
		productGroup.PUT("/:id", handlers.UpdateProduct)
		productGroup.DELETE("/:id", handlers.DeleteProduct)
	}
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handlers.Singup)
		authGroup.POST("/login", handlers.Login)
		authGroup.GET("/validate", middleware.RequireAuth, handlers.Validate)
	}
	orderGroup := router.Group("/order")
	{
		orderGroup.POST("/sell", middleware.RequireAuth, handlers.CreateSellOrder)
		orderGroup.GET("/sell", middleware.RequireAuth, handlers.GetSellOrdersByUID) // ควบรวม
		orderGroup.POST("/buy", middleware.RequireAuth, handlers.CreateBuyOrder)
		orderGroup.GET("/buy", middleware.RequireAuth, handlers.GetBuyOrdersByUID) // ควบรวม
	}
	specsGroup := router.Group("/specs")
	{
		specsGroup.GET("/mouse", middleware.RequireAuth)
		specsGroup.GET("/mousepad", middleware.RequireAuth)
		specsGroup.GET("/cpu", middleware.RequireAuth)
	}
}
