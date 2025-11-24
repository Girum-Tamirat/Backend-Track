package models

// Task represents a task item.
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Status      string `json:"status" binding:"required"` // e.g., "pending", "in-progress", "done"
}
