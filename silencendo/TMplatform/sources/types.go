package sources

import "time"

type KnowledgeSource struct {
	ID        string            `json:"id"`
	Number    int               `json:"number"` // Sequential number for user-friendly access
	Type      string            `json:"type"`   // "file" or "url"
	Path      string            `json:"path"`   // file path or URL
	Title     string            `json:"title"`
	CreatedAt time.Time         `json:"createdAt"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type SourceState struct {
	AllSources       []KnowledgeSource `json:"allSources"`
	ActiveSourceIds  []string          `json:"activeSourceIds"`
	NextSourceNumber int               `json:"nextSourceNumber"` // For sequential numbering
}
