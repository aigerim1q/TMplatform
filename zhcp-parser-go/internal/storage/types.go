package storage

import (
	"context"

	"zhcp-parser-go/internal/transformers"
)

// PersistResult describes what was created in the database.
type PersistResult struct {
	ProjectID int64
	PhaseIDs  []int64
	TaskIDs   []int64
}

// Store persists parsed project structures.
type Store interface {
	Init(ctx context.Context) error
	Close() error
	PersistProjectStructure(ctx context.Context, structure *transformers.ProjectStructure) (*PersistResult, error)
}
