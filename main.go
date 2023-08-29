// main.go
package main

import (
	"ezcom/db"
	"ezcom/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to MongoDB
	db.Connect()

	if err := db.InitFirebaseApp(); err != nil {
		panic("Failed to initialize Firebase: " + err.Error())
	}

	// Set up Gin router
	router := gin.Default()
	routes.Setup(router)

	// Start the server
	err := router.Run(":9090")
	if err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
