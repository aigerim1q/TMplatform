package models

import "time"

// User represents an actor interacting with the chatbot.
type User struct {
	ID        string
	Name      string
	Email     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Project represents a container for stages and tasks.
type Project struct {
	ID               string
	OwnerID          string
	Title            string
	Description      string
	Status           string
	NormalizedTitle  string
	NextTaskStableID int64 // Counter for the next stable ID to assign to new tasks
	DeletedAt        *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Stage represents a project stage.
type Stage struct {
	ID              string
	ProjectID       string
	Title           string
	NormalizedTitle string
	OrderIndex      int
	ResponsibleID   *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Task represents a task optionally bound to a stage.
type Task struct {
	ID              string
	ProjectID       string
	StageID         *string
	Title           string
	NormalizedTitle string
	Description     string
	Status          string
	Priority        string
	AssigneeID      *string
	DueDate         *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	NumericID       int64 // Stable numeric ID for referencing tasks
}

// ProjectDetails aggregates a project with its nested data for read APIs.
type ProjectDetails struct {
	Project Project
	Stages  []Stage
	Tasks   []Task
}
