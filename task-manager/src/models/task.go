package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description"`
	Type         string    `json:"type" gorm:"default:Task;check:type IN ('Story', 'Bug', 'Task', 'Epic')"`
	Status       string    `json:"status" gorm:"default:'To Do'"`
	Priority     string    `json:"priority" gorm:"default:Medium;check:priority IN ('Low', 'Medium', 'High', 'Critical')"`
	ProjectID    string    `json:"project_id" gorm:"not null"`
	ReporterID   string    `json:"reporter_id" gorm:"not null"`
	AssigneeID   *string   `json:"assignee_id"`
	SprintID     *string   `json:"sprint_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DueDate      string    `json:"due_date"`
	TimeEstimate *int      `json:"time_estimate"`
	TimeSpent    *int      `json:"time_spent"`

	// Define relationships
	Project  Project `json:"project" gorm:"foreignKey:ProjectID"`
	Reporter User    `json:"reporter" gorm:"foreignKey:ReporterID"`
	Assignee *User   `json:"assignee" gorm:"foreignKey:AssigneeID"`
}

// BeforeCreate will set default values and a UUID if ID is not set
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	// Set default values if not provided
	if t.Type == "" {
		t.Type = "Task"
	}
	if t.Status == "" {
		t.Status = "To Do"
	}
	if t.Priority == "" {
		t.Priority = "Medium"
	}

	return nil
}

// TasksCount counts the number of tasks for a project with optional status filter
func TasksCount(db *gorm.DB, projectID string, status string) (int64, error) {
	var count int64
	query := db.Model(&Task{}).Where("project_id = ?", projectID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetProjectTaskCounts returns counts of tasks by status for a project
func GetProjectTaskCounts(db *gorm.DB, projectID string) (map[string]int64, error) {
	// Define statuses we want to count
	statuses := []string{"To Do", "In Progress", "Done"}
	counts := make(map[string]int64)

	// Get total count
	total, err := TasksCount(db, projectID, "")
	if err != nil {
		return counts, err
	}
	counts["Total"] = total

	// Get count for each status
	for _, status := range statuses {
		count, err := TasksCount(db, projectID, status)
		if err != nil {
			return counts, err
		}
		counts[status] = count
	}

	return counts, nil
}
