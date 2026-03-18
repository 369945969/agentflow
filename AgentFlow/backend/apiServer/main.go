package main

import (
	"apiServer/config"
	"apiServer/db"
	"apiServer/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (for secrets or traditional env vars)
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Load configuration file
	config.LoadConfig()

	// Initialize database
	db.InitDB()
	defer db.DB.Close()

	// Sync skills from configured directory
	db.SyncSkills()

	// Initialize router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.Default())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Backend server is running (Go)",
		})
	})

	// Register routes
	routes.RegisterSkillRoutes(r)
	routes.RegisterAgentRoutes(r)
	routes.RegisterModelRoutes(r)

	// Start server
	port := config.AppConfig.Server.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server is listening at http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
