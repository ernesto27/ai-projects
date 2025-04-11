package api

import (
	"net/http"
	"os"

	"github.com/ernesto/task-manager/src/auth"
	"github.com/ernesto/task-manager/src/config"
	"github.com/ernesto/task-manager/src/models"
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
	router.LoadHTMLGlob("templates/*.tmpl")
	// Serve static files
	router.Static("/static", "./static")

	// Web routes (HTML)
	router.GET("/", HomeHandler)
	router.GET("/login", LoginHandler)
	router.GET("/register", RegisterHandler)
	router.GET("/dashboard", auth.AuthMiddlewareWeb(), DashboardHandler)
	router.GET("/projects", auth.AuthMiddlewareWeb(), ProjectsPageHandler)
	router.GET("/projects/new", auth.AuthMiddlewareWeb(), NewProjectPageHandler)
	router.GET("/projects/:id", auth.AuthMiddlewareWeb(), ProjectDetailsPageHandler)
	router.GET("/projects/:id/edit", auth.AuthMiddlewareWeb(), EditProjectPageHandler)

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

		// Project routes
		protected.GET("/projects", GetProjects)
		protected.POST("/projects", CreateProject)
		protected.GET("/projects/:id", GetProject)
		protected.PUT("/projects/:id", UpdateProject)
		protected.DELETE("/projects/:id", DeleteProject)
	}

	return router
}

// HomeHandler renders the home page
func HomeHandler(c *gin.Context) {
	c.HTML(200, "index", gin.H{
		"Title": "Home",
	})
}

// RegisterHandler renders the registration page
func RegisterHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", gin.H{
		"Title": "Register",
	})
}

// LoginHandler renders the login page
func LoginHandler(c *gin.Context) {
	c.HTML(200, "login.tmpl", gin.H{
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

	c.HTML(200, "dashboard.tmpl", gin.H{
		"Title": "Dashboard",
		"User":  user,
	})
}

// ProjectsPageHandler renders the projects list page
func ProjectsPageHandler(c *gin.Context) {
	var projects []models.Project

	// Get projects from database with owner information
	if err := config.DB.Preload("Owner").Find(&projects).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Failed to fetch projects",
		})
		return
	}

	// Get current user
	user, _ := c.Get("user")

	c.HTML(http.StatusOK, "projects.tmpl", gin.H{
		"Title":    "Projects",
		"User":     user,
		"Projects": projects,
	})
}

// NewProjectPageHandler renders the new project form
func NewProjectPageHandler(c *gin.Context) {
	// Get all users to populate owner selection dropdown
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Failed to fetch users",
		})
		return
	}

	// Get current user
	currentUser, _ := c.Get("user")

	c.HTML(http.StatusOK, "project_form.tmpl", gin.H{
		"Title":      "Create New Project",
		"User":       currentUser,
		"Users":      users,
		"FormAction": "/api/projects",
		"Method":     "POST",
		"Project":    models.Project{}, // Empty project for new form
	})
}

// ProjectDetailsPageHandler renders the project details page
func ProjectDetailsPageHandler(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	// Get project from database with owner information
	if err := config.DB.Preload("Owner").First(&project, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Project not found",
		})
		return
	}

	// Get current user
	user, _ := c.Get("user")

	c.HTML(http.StatusOK, "project_details.tmpl", gin.H{
		"Title":   project.Name,
		"User":    user,
		"Project": project,
	})
}

// EditProjectPageHandler renders the edit project form
func EditProjectPageHandler(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	// Get project from database
	if err := config.DB.First(&project, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Project not found",
		})
		return
	}

	// Get all users to populate owner selection dropdown
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Failed to fetch users",
		})
		return
	}

	// Get current user
	currentUser, _ := c.Get("user")

	c.HTML(http.StatusOK, "project_form.tmpl", gin.H{
		"Title":      "Edit Project",
		"User":       currentUser,
		"Users":      users,
		"FormAction": "/api/projects/" + id,
		"Method":     "PUT",
		"Project":    project,
	})
}
