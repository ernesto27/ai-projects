package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"html/template"

	"github.com/ernesto/task-manager/src/auth"
	"github.com/ernesto/task-manager/src/config"
	"github.com/ernesto/task-manager/src/models"
	"github.com/ernesto/task-manager/src/utils"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter() *gin.Engine {
	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Set up HTML templates with custom functions
	tmpl := template.Must(template.New("").Funcs(utils.GetTemplateFunctions()).ParseGlob("templates/*.tmpl"))
	router.SetHTMLTemplate(tmpl)

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
	router.GET("/projects/:id/tasks", auth.AuthMiddlewareWeb(), ProjectTasksPageHandler)
	router.GET("/projects/:id/tasks/new", auth.AuthMiddlewareWeb(), NewTaskPageHandler)           // New route for task form
	router.GET("/projects/:id/tasks/:taskId", auth.AuthMiddlewareWeb(), TaskDetailsPageHandler)   // Task details page
	router.GET("/projects/:id/tasks/:taskId/edit", auth.AuthMiddlewareWeb(), EditTaskPageHandler) // Edit task form

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

		// Task routes
		protected.GET("/projects/:id/tasks", GetTasksHandler)
		protected.POST("/projects/:id/tasks", CreateTaskHandler)
		protected.GET("/projects/:id/tasks/:taskId", GetTaskHandler)
		protected.PUT("/projects/:id/tasks/:taskId", UpdateTaskHandler)
		protected.DELETE("/projects/:id/tasks/:taskId", DeleteTaskHandler)
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

	// Fetch projects from the database
	var projects []models.Project
	if err := config.DB.Preload("Owner").Limit(5).Order("created_at DESC").Find(&projects).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Failed to fetch projects",
		})
		return
	}

	// Count total projects
	var projectCount int64
	config.DB.Model(&models.Project{}).Count(&projectCount)

	// Count total users
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)

	c.HTML(200, "dashboard.tmpl", gin.H{
		"Title":        "Dashboard",
		"User":         user,
		"Projects":     projects,
		"ProjectCount": projectCount,
		"UserCount":    userCount,
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

// ProjectTasksPageHandler renders the tasks page for a specific project
func ProjectTasksPageHandler(c *gin.Context) {
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

	// Get tasks for this project with all related data
	var tasks []models.Task
	if err := config.DB.Where("project_id = ?", id).
		Preload("Assignee").
		Preload("Reporter").
		Order("id DESC").
		Find(&tasks).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Failed to fetch tasks",
		})
		return
	}

	// Get all users for the assignee dropdown
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Failed to fetch users",
		})
		return
	}

	// Get current user
	user, _ := c.Get("user")

	c.HTML(http.StatusOK, "tasks.tmpl", gin.H{
		"Title":   project.Name + " - Tasks",
		"User":    user,
		"Project": project,
		"Tasks":   tasks,
		"Users":   users,
	})
}

// NewTaskPageHandler renders the new task form
func NewTaskPageHandler(c *gin.Context) {
	projectID := c.Param("id")
	var project models.Project

	// Get project from database
	if err := config.DB.First(&project, projectID).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Project not found",
		})
		return
	}

	// Get all users to populate assignee selection dropdown
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

	c.HTML(http.StatusOK, "task_form.tmpl", gin.H{
		"Title":      "Create New Task",
		"User":       currentUser,
		"Users":      users,
		"FormAction": "/api/projects/" + projectID + "/tasks",
		"ProjectID":  projectID,
		"Method":     "POST",
		"Task":       models.Task{}, // Empty task for new form
	})
}

// TaskDetailsPageHandler renders the task details page
func TaskDetailsPageHandler(c *gin.Context) {
	projectID := c.Param("id")
	taskID := c.Param("taskId")
	var project models.Project
	var task models.Task

	// Get project from database
	if err := config.DB.First(&project, projectID).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Project not found",
		})
		return
	}

	// Get task with related entities
	if err := config.DB.Where("id = ? AND project_id = ?", taskID, projectID).
		Preload("Reporter").
		Preload("Assignee").
		First(&task).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Task not found",
		})
		return
	}

	// Get current user
	currentUser, _ := c.Get("user")

	c.HTML(http.StatusOK, "task_details.tmpl", gin.H{
		"Title":       task.Title,
		"User":        currentUser,
		"Project":     project,
		"Task":        task,
		"CurrentUser": currentUser,
	})
}

// EditTaskPageHandler renders the edit task form
func EditTaskPageHandler(c *gin.Context) {
	projectID := c.Param("id")
	taskID := c.Param("taskId")
	var project models.Project
	var task models.Task

	// Get project from database
	if err := config.DB.First(&project, projectID).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Project not found",
		})
		return
	}

	// Get task with related entities
	if err := config.DB.Where("id = ? AND project_id = ?", taskID, projectID).First(&task).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Task not found",
		})
		return
	}

	t, err := time.Parse(time.RFC3339, task.DueDate)
	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Task not found",
		})
		return
	}

	output := t.Format("2006-01-02")
	task.DueDate = output

	// Get all users to populate assignee selection dropdown
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

	c.HTML(http.StatusOK, "task_form.tmpl", gin.H{
		"Title":      "Edit Task",
		"User":       currentUser,
		"Users":      users,
		"FormAction": "/api/projects/" + projectID + "/tasks/" + taskID,
		"ProjectID":  projectID,
		"Method":     "PUT",
		"Task":       task,
	})
}

// GetTasks handles fetching tasks for a project
func GetTasksHandler(c *gin.Context) {
	projectID := c.Param("id")
	var tasks []models.Task

	if err := config.DB.Where("project_id = ?", projectID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// CreateTaskHandler handles creating a new task for a project
func CreateTaskHandler(c *gin.Context) {
	projectID := c.Param("id")
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Get current user from context
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user := currentUser.(models.User)

	// Set task properties
	task.ProjectID = projectID
	task.ReporterID = strconv.FormatUint(uint64(user.ID), 10) // Convert uint ID to string
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Set default values if not provided
	if task.Status == "" {
		task.Status = "To Do"
	}
	if task.Type == "" {
		task.Type = "Task"
	}
	if task.Priority == "" {
		task.Priority = "Medium"
	}

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTask handles fetching a specific task
func GetTaskHandler(c *gin.Context) {
	taskID := c.Param("taskId")
	var task models.Task

	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask handles updating a specific task
func UpdateTaskHandler(c *gin.Context) {
	taskID := c.Param("taskId")
	var task models.Task

	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask handles deleting a specific task
func DeleteTaskHandler(c *gin.Context) {
	taskID := c.Param("taskId")
	var task models.Task

	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
