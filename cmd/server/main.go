package main

import (
	"log"
	"os"

	"civra_back/internal/database"
	"civra_back/internal/handlers"
	"civra_back/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	log.Println("🚀 Initializing database connection...")
	if err := database.InitDatabase(); err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	log.Println("✅ Database connected successfully!")

	// Graceful shutdown
	defer database.CloseDatabase()

	// Create Gin router
	r := gin.Default()

	// --- CORS Middleware ---
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// --- Public routes ---
	userRoutes := r.Group("/api/users")
	{
		userRoutes.POST("/register", handlers.CreateUser)
		userRoutes.POST("/login", handlers.LoginUser)
		userRoutes.GET("/:id", handlers.GetUser)
		userRoutes.GET("/", handlers.GetUsers) // for testing
	}

	// --- Protected routes ---
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", handlers.GetProfile) // example protected route
	}

	// --- Health check ---
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"message":  "Server is running",
			"database": "Connected",
		})
	})

	// --- Run server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🌍 Server is running on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
