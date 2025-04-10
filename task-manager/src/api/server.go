package api

import (
	"html/template"
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

	// Setup template rendering with proper inheritance
	templ := template.Must(template.ParseFiles(
		"templates/layouts/base.html",
		"templates/index.html",
		"templates/dashboard.html",
		"templates/auth/login.html",
		"templates/auth/register.html",
	))
	router.SetHTMLTemplate(templ)

	// Serve static files
	router.Static("/static", "./static")

	// Web routes (HTML)
	router.GET("/", HomeHandler)
	router.GET("/register", RegisterHandler)
	router.GET("/login", LoginHandler)
	router.GET("/dashboard", auth.AuthMiddlewareWeb(), DashboardHandler)

	// Public API routes (no authentication required)
	public := router.Group("/api")
	{
		// Authentication routes
		public.POST("/register", auth.Register) // Email registration
		public.POST("/login", auth.Login)       // Email login
	}

	// Protected API routes (authentication required)
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

// HomeHandler renders the home page
func HomeHandler(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"Title": "Home",
	})
}

// RegisterHandler renders the registration page
func RegisterHandler(c *gin.Context) {
	c.HTML(200, "auth/register.html", gin.H{
		"Title": "Register",
	})
}

// LoginHandler renders the login page
func LoginHandler(c *gin.Context) {
	c.HTML(200, "auth/login.html", gin.H{
		"Title": "Login",
	})
}

// DashboardHandler renders the dashboard page (protected)
func DashboardHandler(c *gin.Context) {
	// Get current user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		c.Redirect(302, "/login")
		return
	}

	c.HTML(200, "dashboard.html", gin.H{
		"Title": "Dashboard",
		"User":  user,
	})
}
