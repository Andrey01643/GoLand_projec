package models

import "time"

// Task - модель задачи
type Task struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      int       `json:"user_id"`
}
