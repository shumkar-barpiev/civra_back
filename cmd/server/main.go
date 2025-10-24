package main

import (
	"log"
	"os"

	"civra_back/internal/database"
	"civra_back/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	log.Println("🚀 Initializing database connection...")
	err := database.InitDatabase()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	log.Println("✅ Database connected successfully!")

	// Setup graceful shutdown
	defer database.CloseDatabase()

	// Create router
	r := gin.Default()

	// Add CORS middleware (important for frontend communication)
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

	// User routes
	userRoutes := r.Group("/api/users")
	{
		userRoutes.POST("/register", handlers.CreateUser)
		userRoutes.POST("/login", handlers.LoginUser)
		userRoutes.GET("/:id", handlers.GetUser)
		userRoutes.GET("/", handlers.GetUsers) // For testing
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
			"database": "Connected",
		})
	})

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("🌐 Server starting on :%s", port)
	log.Printf("📚 API Documentation:")
	log.Printf("   POST /api/users/register - Register a new user")
	log.Printf("   POST /api/users/login    - User login")
	log.Printf("   GET  /api/users/:id      - Get user by ID")
	log.Printf("   GET  /health             - Health check")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}