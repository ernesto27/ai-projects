package api

import (
	"net/http"
	"strconv"

	"github.com/ernesto/task-manager/src/config"
	"github.com/ernesto/task-manager/src/models"
	"github.com/ernesto/task-manager/src/utils"
	"github.com/gin-gonic/gin"
)

// CreateProject creates a new project
func CreateProject(c *gin.Context) {
	var input models.ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	// Check if project key already exists
	var existingProject models.Project
	if result := config.DB.Where("key = ?", input.Key).First(&existingProject); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project key already exists"})
		return
	}

	// Get current user from context (set by auth middleware)
	_, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// No need to assign to a variable if we're not using it
	// Just check that it exists for authentication

	// Create project
	project := models.Project{
		Name:        input.Name,
		Key:         input.Key,
		Description: input.Description,
		OwnerID:     input.OwnerID,
	}

	// Save project to database
	result := config.DB.Create(&project)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Project created successfully",
		"data":    project,
	})
}

// GetProjects returns all projects
func GetProjects(c *gin.Context) {
	var projects []models.Project

	// Preload owner information
	result := config.DB.Preload("Owner").Find(&projects)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(projects),
		"data":   projects,
	})
}

// GetProject returns a specific project by ID
func GetProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	// Preload owner information
	result := config.DB.Preload("Owner").First(&project, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   project,
	})
}

// UpdateProject updates a project
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	// Find project by ID
	result := config.DB.First(&project, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Parse input
	var input models.ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	// Check if key is being changed and if it already exists
	if project.Key != input.Key {
		var existingProject models.Project
		if result := config.DB.Where("key = ? AND id != ?", input.Key, id).First(&existingProject); result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project key already exists"})
			return
		}
	}

	// Update project fields
	project.Name = input.Name
	project.Key = input.Key
	project.Description = input.Description
	project.OwnerID = input.OwnerID

	// Save updated project
	result = config.DB.Save(&project)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Project updated successfully",
		"data":    project,
	})
}

// DeleteProject deletes a project
func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	// Convert ID to uint
	projectID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Delete project
	result := config.DB.Delete(&models.Project{}, uint(projectID))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	// Check if no rows were affected (project wasn't found)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Project deleted successfully",
	})
}
