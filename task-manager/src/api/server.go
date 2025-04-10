package api

import (
	"os"

	"github.com/ernesto/task-manager/src/auth"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter() *gin.Engine {
	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Public routes (no authentication required)
	public := router.Group("/api")
	{
		// Authentication routes
		public.POST("/register", auth.Register) // Email registration
		public.POST("/login", auth.Login)       // Email login
	}

	// Protected routes (authentication required)
	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware())
	{
		// User profile routes
		protected.GET("/profile", GetUserProfile)
		protected.PUT("/profile", UpdateUserProfile)
		protected.POST("/profile/change-password", ChangePassword)
	}

	return router
}
