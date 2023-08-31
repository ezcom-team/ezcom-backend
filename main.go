// main.go
package main

import (
	"ezcom/db"
	"ezcom/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
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

	// Set up CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // แก้ไข URL ของโดเมน React ของคุณตรงนี้
	router.Use(cors.New(config))

	routes.Setup(router)

	// Start the server
	var port = envPortOr("3000")
	err := router.Run(port)
	if err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}

func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
}
