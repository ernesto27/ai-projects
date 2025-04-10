package api

import (
	"github.com/ernesto/task-manager/src/auth"
	"github.com/ernesto/task-manager/src/config"
	"github.com/ernesto/task-manager/src/utils"
	"github.com/gin-gonic/gin"
)

// ProfileUpdateInput represents the data needed for updating a user's profile
type ProfileUpdateInput struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// GetUserProfile returns the authenticated user's profile
func GetUserProfile(c *gin.Context) {
	// Get current user from auth middleware
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// Return user profile data
	c.JSON(200, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"role":      user.Role,
			"avatar":    user.Avatar,
			"createdAt": user.CreatedAt,
		},
	})
}

// UpdateUserProfile updates the authenticated user's profile information
func UpdateUserProfile(c *gin.Context) {
	// Get current user from auth middleware
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// Bind input data
	var input ProfileUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.Name != "" {
		user.Name = utils.SanitizeInput(input.Name)
	}

	if input.Avatar != "" {
		user.Avatar = utils.SanitizeInput(input.Avatar)
	}

	// Save changes to database
	if err := config.DB.Save(user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update profile"})
		return
	}

	// Return updated user profile
	c.JSON(200, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"role":      user.Role,
			"avatar":    user.Avatar,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

// ChangePassword allows users to update their password
func ChangePassword(c *gin.Context) {
	// Get current user from auth middleware
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// Bind input data
	var input struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verify current password
	if err := user.CheckPassword(input.CurrentPassword); err != nil {
		c.JSON(401, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Validate new password
	if valid, message := utils.ValidatePassword(input.NewPassword); !valid {
		c.JSON(400, gin.H{"error": message})
		return
	}

	// Hash the new password
	if err := user.HashPassword(input.NewPassword); err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// Save changes to database
	if err := config.DB.Save(user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(200, gin.H{"message": "Password changed successfully"})
}
