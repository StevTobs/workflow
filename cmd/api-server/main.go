package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"

	auth "item-workflow-system/internal/auth"
	metrics "item-workflow-system/internal/metrics"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func handler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	// your existing code here
	main()

}
func main() {
	// Initialize the DB connection
	metrics.InitDB()
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Adjust as needed
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Authentication routes
	r.POST("/login", auth.LoginHandler)
	r.POST("/logout", auth.LogoutHandler)
	r.GET("/test/items", metrics.ReadAllItems_)

	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/items", metrics.CreateItem_)
		protected.GET("/items/:id", metrics.ReadItemByID)
		protected.GET("/items", metrics.ReadAllItems_)
		protected.PUT("/items/:id", metrics.UpdateItem_)
		protected.PATCH("/items/:id/status", metrics.PatchItemStatusHandler)
		protected.DELETE("/items/:id", metrics.RemoveItem_)
	}

	// Create a server
	srv := &http.Server{
		Addr:    ":3456",
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

	// Allow CORS for all origins (use with caution in production)
	corsObj := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":3456", handlers.CORS(corsObj)(http.DefaultServeMux))

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		gin.DefaultWriter.Write([]byte("Server shutdown failed: " + err.Error()))
	}

	gin.DefaultWriter.Write([]byte("Server exited properly"))
}
