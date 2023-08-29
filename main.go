// main.go
package main

import (
	"ezcom/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to MongoDB

	// Set up Gin router
	router := gin.Default()
	routes.Setup(router)

	// Start the server
	// Returns PORT from environment if found, defaults to
	// value in `port` parameter otherwise. The returned port
	// is prefixed with a `:`, e.g. `":3000"`.

	var port = envPortOr("0303")
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
