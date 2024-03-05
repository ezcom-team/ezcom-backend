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
		productGroup.GET("/spec/:type/:id", handlers.GetSpecByID)
		productGroup.GET("/specs/:type", handlers.GetSpecs)
		productGroup.PUT("/:id", handlers.UpdateProduct)
		productGroup.DELETE("/:id", handlers.DeleteProduct)
	}
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handlers.Singup)
		authGroup.POST("/login", handlers.Login)
		authGroup.GET("/validate", middleware.RequireAuth, handlers.Validate)
	}
	userGroup := router.Group("/user")
	{
		userGroup.GET("/user", middleware.RequireAuth, handlers.GetUserProfile)
		userGroup.GET("/users", middleware.RequireAuth, handlers.GetUsers)
		userGroup.PUT("/:id", middleware.RequireAuth, handlers.UpdateUser)
		userGroup.DELETE("/:id", middleware.RequireAuth, handlers.DeleteUser)
	}
	orderGroup := router.Group("/order")
	{
		orderGroup.POST("/sell", middleware.RequireAuth, handlers.CreateSellOrder)
		orderGroup.GET("/sell", middleware.RequireAuth, handlers.GetSellOrdersByUID)       
		orderGroup.GET("/sells", middleware.RequireAuth, handlers.GetSellOrders)           
		orderGroup.GET("/sells/:pid", middleware.RequireAuth, handlers.GetSellOrdersByPID) 
		orderGroup.DELETE("/:type/:id", middleware.RequireAuth, handlers.DeleteOrder)      
		orderGroup.POST("/buy", middleware.RequireAuth, handlers.CreateBuyOrder)
		orderGroup.GET("/buy", middleware.RequireAuth, handlers.GetBuyOrdersByUID)                                         
		orderGroup.GET("/buys", middleware.RequireAuth, handlers.GetBuyOrders)                                            
		orderGroup.GET("/buys/:pid", middleware.RequireAuth, handlers.GetBuyOrdersByPID)                                  
		orderGroup.GET("/:type/:id", middleware.RequireAuth, handlers.DeleteOrder)                                        
		orderGroup.GET("/matched-order", middleware.RequireAuth, handlers.GetMatchedOrder)                                 
		orderGroup.GET("/matched-order/:pid", middleware.RequireAuth, handlers.GetMatchedOrderByPID)                       
		orderGroup.PUT("/matched-order/status", middleware.RequireAuth, handlers.UpdataMatchedOrderStatus)                 
		orderGroup.PUT("/matched-order/received", middleware.RequireAuth, handlers.UpdataMatchedOrderRecived)              
		orderGroup.PUT("/matched-order/tracking-number", middleware.RequireAuth, handlers.UpdataMatchedOrderTackingNumber) 
		orderGroup.PUT("/matched-order/payment-status", middleware.RequireAuth, handlers.UpdataPaymentStatus)
	}
	specsGroup := router.Group("/specs")
	{
		specsGroup.GET("/mouse", middleware.RequireAuth)
		specsGroup.GET("/mousepad", middleware.RequireAuth)
		specsGroup.GET("/cpu", middleware.RequireAuth)
	}
}
