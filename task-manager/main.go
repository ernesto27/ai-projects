package main

import (
	"fmt"
	"os"

	"github.com/ernesto/task-manager/src/api"
	"github.com/ernesto/task-manager/src/config"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: Error loading .env file. Using environment variables.")
	}

	// Connect to database
	config.ConnectDatabase()

	// Set up router
	router := api.SetupRouter()

	// Get port from environment or use default
	port := getEnv("PORT", "8080")

	// Start server
	fmt.Printf("Server running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
