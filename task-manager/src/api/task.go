package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ernesto/task-manager/src/config"
	"github.com/ernesto/task-manager/src/models"
)

// GetTasks handles fetching tasks for a project
func GetTasks(c *gin.Context) {
	projectID := c.Param("id")
	var tasks []models.Task

	// Get tasks with related entities
	if err := config.DB.Where("project_id = ?", projectID).
		Preload("Reporter").
		Preload("Assignee").
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// CreateTask handles creating a new task for a project
func CreateTask(c *gin.Context) {
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

	// Create the task
	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task: " + err.Error()})
		return
	}

	// Return the created task
	c.JSON(http.StatusCreated, task)
}

// GetTask handles fetching a specific task
func GetTask(c *gin.Context) {
	taskID := c.Param("taskId")
	var task models.Task

	// Get task with related entities
	if err := config.DB.
		Preload("Reporter").
		Preload("Assignee").
		First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask handles updating a specific task
func UpdateTask(c *gin.Context) {
	taskID := c.Param("taskId")
	var task models.Task

	// Find the existing task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Store original values that shouldn't change
	createdAt := task.CreatedAt
	reporterID := task.ReporterID

	// Bind JSON data to task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Preserve fields that shouldn't be changed

	task.CreatedAt = createdAt
	task.ReporterID = reporterID
	task.UpdatedAt = time.Now()

	// Update the task
	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask handles deleting a specific task
func DeleteTask(c *gin.Context) {
	taskID := c.Param("taskId")
	var task models.Task

	// Find the task to delete
	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Delete the task
	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// TasksPageHandler renders the tasks page for a project
func TasksPageHandler(c *gin.Context) {
	projectID := c.Param("id")
	var project models.Project
	var tasks []models.Task
	var users []models.User

	// Get project details
	if err := config.DB.First(&project, projectID).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Project not found",
		})
		return
	}

	// Get tasks for the project with related entities
	if err := config.DB.Where("project_id = ?", projectID).
		Preload("Reporter").
		Preload("Assignee").
		Find(&tasks).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Failed to fetch tasks",
		})
		return
	}

	// Get all users for assignee selection
	if err := config.DB.Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"Title":   "Error",
			"Message": "Failed to fetch users",
		})
		return
	}

	// Get current user
	currentUser, _ := c.Get("user")

	c.HTML(http.StatusOK, "tasks.tmpl", gin.H{
		"Title":       "Tasks - " + project.Name,
		"User":        currentUser,
		"Project":     project,
		"Tasks":       tasks,
		"Users":       users,
		"CurrentUser": currentUser,
	})
}
