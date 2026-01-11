package cli

import "time"

// ConfirmationKind enumerates pending confirmation types.
type ConfirmationKind string

const (
	ConfirmationProjectVsDocument ConfirmationKind = "project_vs_document"
)

// PendingConfirmation holds in-flight disambiguation.
type PendingConfirmation struct {
	Kind          ConfirmationKind
	OriginalInput string
	CreatedAt     time.Time
}

// PendingDeletion captures confirmation required before soft-deleting a project.
type PendingDeletion struct {
	ProjectID    string
	ProjectTitle string
}

// SessionState stores lightweight, in-memory session data for the CLI.
type SessionState struct {
	Pending            *PendingConfirmation
	ActiveProjectID    string
	ActiveProjectTitle string
	LastProjectList    []ProjectSummary
	PendingDeletion    *PendingDeletion
}

func NewSessionState() *SessionState {
	return &SessionState{LastProjectList: []ProjectSummary{}}
}

// ProjectSummary caches the last shown list for numeric selection.
type ProjectSummary struct {
	ID              string
	Title           string
	Status          string
	NormalizedTitle string
}
