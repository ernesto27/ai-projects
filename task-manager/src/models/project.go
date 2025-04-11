package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Project represents a project in the system
type Project struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Key         string    `json:"key" gorm:"uniqueIndex;size:10;not null"` // Short code like "TASK" or "JIRA"
	Description string    `json:"description"`
	OwnerID     uint      `json:"ownerId" gorm:"not null"`
	Owner       User      `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ProjectInput represents the data needed for creating/updating a project
type ProjectInput struct {
	Name        string `json:"name" binding:"required"`
	Key         string `json:"key" binding:"required,min=2,max=10,alphanum"`
	Description string `json:"description"`
	OwnerID     uint   `json:"ownerId" binding:"required"`
}

// BeforeCreate hook is called before creating a Project instance
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	// Convert key to uppercase
	p.Key = stringToUpper(p.Key)
	return nil
}

// BeforeUpdate hook is called before updating a Project instance
func (p *Project) BeforeUpdate(tx *gorm.DB) error {
	// Convert key to uppercase
	p.Key = stringToUpper(p.Key)
	return nil
}

// Helper function to convert string to uppercase
func stringToUpper(s string) string {
	if s != "" {
		return strings.ToUpper(s)
	}
	return ""
}
