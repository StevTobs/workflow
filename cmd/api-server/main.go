package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	auth "item-workflow-system/internal/auth"
	metrics "item-workflow-system/internal/metrics"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the DB connection
	metrics.InitDB()
	r := gin.Default()

	// Authentication routes
	r.POST("/login", auth.LoginHandler)
	r.POST("/logout", auth.LogoutHandler)

	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/items", metrics.CreateItem_)
		protected.GET("/items", metrics.ReadAllItems_)
		protected.PUT("/items/:id", metrics.UpdateItem_)
		protected.PATCH("/items/:id/status", metrics.PatchItemStatusHandler)
		protected.DELETE("/items/:id", metrics.RemoveItem_)
	}

	// Create a server
	srv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	// Run the server in a goroutine so that it doesn't block
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			gin.DefaultWriter.Write([]byte("Server failed to start: " + err.Error()))
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Wait for interrupt signal

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		gin.DefaultWriter.Write([]byte("Server shutdown failed: " + err.Error()))
	}

	gin.DefaultWriter.Write([]byte("Server exited properly"))
}
