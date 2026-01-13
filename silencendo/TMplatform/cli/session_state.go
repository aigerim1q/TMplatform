package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
	state := &SessionState{LastProjectList: []ProjectSummary{}}
	state.loadFromDisk()
	return state
}

const sessionStateFile = ".bot/session_state.json"

// ProjectSummary caches the last shown list for numeric selection.
type ProjectSummary struct {
	ID              string
	Title           string
	Status          string
	NormalizedTitle string
}

func (s *SessionState) SetActiveProject(id, title string) {
	if s == nil {
		return
	}
	s.ActiveProjectID = strings.TrimSpace(id)
	s.ActiveProjectTitle = strings.TrimSpace(title)
	s.persist()
}

func (s *SessionState) ClearActiveProject() {
	if s == nil {
		return
	}
	s.ActiveProjectID = ""
	s.ActiveProjectTitle = ""
	s.persist()
}

func (s *SessionState) loadFromDisk() {
	if s == nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(sessionStateFile), 0o755); err != nil {
		return
	}
	data, err := os.ReadFile(sessionStateFile)
	if err != nil {
		return
	}
	var disk struct {
		ActiveProjectID    string `json:"active_project_id"`
		ActiveProjectTitle string `json:"active_project_title"`
	}
	if err := json.Unmarshal(data, &disk); err != nil {
		return
	}
	s.ActiveProjectID = strings.TrimSpace(disk.ActiveProjectID)
	s.ActiveProjectTitle = strings.TrimSpace(disk.ActiveProjectTitle)
}

func (s *SessionState) persist() {
	if s == nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(sessionStateFile), 0o755); err != nil {
		return
	}
	disk := struct {
		ActiveProjectID    string `json:"active_project_id"`
		ActiveProjectTitle string `json:"active_project_title"`
	}{ActiveProjectID: strings.TrimSpace(s.ActiveProjectID), ActiveProjectTitle: strings.TrimSpace(s.ActiveProjectTitle)}
	data, err := json.MarshalIndent(disk, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(sessionStateFile, data, 0o644)
}
