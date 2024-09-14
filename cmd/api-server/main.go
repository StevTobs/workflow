package main

import (
	aunt "item-workflow-system/internal/aunt"
	metrics "item-workflow-system/internal/metrics"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the DB connection
	metrics.InitDB()
	r := gin.Default()

	// Authentication routes
	r.POST("/login", aunt.LoginHandler)
	r.POST("/logout", aunt.LogoutHandler)

	// Apply AuthMiddleware to protected routes
	protected := r.Group("/")
	protected.Use(aunt.AuthMiddleware())
	{
		// Define the POST endpoint to create a new item (protected)
		protected.POST("/items", aunt.AuthMiddleware(), metrics.CreateItem)

		protected.GET("/items", aunt.AuthMiddleware(), metrics.ReadAllItems)
		// Define the PUT endpoint to update an item by ID (protected)
		protected.PUT("/items/:id", aunt.AuthMiddleware(), metrics.UpdateItem)
		// Define the PATCH endpoint to update the status of an item by ID (protected)
		protected.PATCH("/items/:id/status", aunt.AuthMiddleware(), metrics.PatialUpdateItem)
		// Define the DELETE endpoint to delete an item by ID (protected)
		protected.DELETE("/items/:id", aunt.AuthMiddleware(), metrics.RemoveItem)

	}

	// Run the server
	r.Run(":2024")
}
