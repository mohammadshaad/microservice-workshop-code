package main

import (
	"log"

	"github.com/mohammadshaad/task-service/db"
	"github.com/mohammadshaad/task-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize MongoDB connection
	db.Connect()

	// Create a new Gin router
	router := gin.Default()

	// Setup CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupTaskRoutes(router)

	// Start the server
	defer db.Disconnect()
	log.Fatal(router.Run(":8081"))
}
