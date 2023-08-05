// routes/routes.go
package routes

import (
	"ezcom/handlers"

	"github.com/gin-gonic/gin"
)

// Setup initializes the routes and handlers
func Setup(router *gin.Engine) {
	productGroup := router.Group("/products")
	{
		productGroup.POST("", handlers.CreateProduct)
		productGroup.GET("/:id", handlers.GetProductByID)
		productGroup.GET("/", handlers.GetProducts)
		productGroup.PUT("/:id", handlers.UpdateProduct)
		productGroup.DELETE("/:id", handlers.DeleteProduct)
	}
}
