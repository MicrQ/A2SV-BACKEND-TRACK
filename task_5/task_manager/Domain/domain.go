package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a task entity with business rules.
type Task struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Status      string    `json:"status" bson:"status"`
}

// Validate checks if the task is valid according to business rules.
func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	if t.Status != "pending" && t.Status != "in_progress" && t.Status != "completed" {
		return errors.New("invalid status")
	}
	return nil
}

// User represents a user entity with business rules.
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	PasswordHash string             `bson:"-"`
	Role         string             `bson:"role" json:"role"` // "admin" or "user"
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// Validate checks if the user is valid.
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Role != "admin" && u.Role != "user" {
		return errors.New("invalid role")
	}
	return nil
}

// IsAdmin checks if the user has admin role.
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
