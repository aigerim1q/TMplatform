package sources

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const SOURCES_FILE = ".bot/sources.json"

type SourceManager struct {
	state SourceState
}

func NewSourceManager() *SourceManager {
	sm := &SourceManager{
		state: SourceState{
			AllSources:       []KnowledgeSource{},
			ActiveSourceIds:  []string{},
			NextSourceNumber: 1,
		},
	}
	sm.loadSources()
	return sm
}

func (sm *SourceManager) ensureBotDir() {
	if _, err := os.Stat(".bot"); os.IsNotExist(err) {
		os.MkdirAll(".bot", 0755)
	}
}

func (sm *SourceManager) loadSources() {
	sm.ensureBotDir()

	if _, err := os.Stat(SOURCES_FILE); err == nil {
		data, err := ioutil.ReadFile(SOURCES_FILE)
		if err == nil {
			var tempState struct {
				AllSources []struct {
					ID        string            `json:"id"`
					Number    int               `json:"number"`
					Type      string            `json:"type"`
					Path      string            `json:"path"`
					Title     string            `json:"title"`
					CreatedAt string            `json:"createdAt"`
					Metadata  map[string]string `json:"metadata,omitempty"`
				} `json:"allSources"`
				ActiveSourceIds  []string `json:"activeSourceIds"`
				NextSourceNumber int      `json:"nextSourceNumber"`
			}

			err = json.Unmarshal(data, &tempState)
			if err == nil {
				// Convert date strings back to time.Time objects
				var loadedSources []KnowledgeSource
				for _, source := range tempState.AllSources {
					createdAt, _ := time.Parse(time.RFC3339, source.CreatedAt)
					loadedSources = append(loadedSources, KnowledgeSource{
						ID:        source.ID,
						Number:    source.Number,
						Type:      source.Type,
						Path:      source.Path,
						Title:     source.Title,
						CreatedAt: createdAt,
						Metadata:  source.Metadata,
					})
				}

				// Calculate nextSourceNumber based on highest number in loaded sources
				maxNumber := 0
				for _, s := range loadedSources {
					if s.Number > maxNumber {
						maxNumber = s.Number
					}
				}
				nextNumber := maxNumber + 1
				if nextNumber < tempState.NextSourceNumber {
					nextNumber = tempState.NextSourceNumber
				}

				sm.state = SourceState{
					AllSources:       loadedSources,
					ActiveSourceIds:  tempState.ActiveSourceIds,
					NextSourceNumber: nextNumber,
				}
			} else {
				// Use default state if loading fails
				sm.saveSources()
			}
		} else {
			// Use default state if loading fails
			sm.saveSources()
		}
	} else {
		// Initialize with empty state
		sm.saveSources()
	}
}

func (sm *SourceManager) saveSources() {
	sm.ensureBotDir()

	// Convert time.Time objects to strings for JSON serialization
	var serializableSources []struct {
		ID        string            `json:"id"`
		Number    int               `json:"number"`
		Type      string            `json:"type"`
		Path      string            `json:"path"`
		Title     string            `json:"title"`
		CreatedAt string            `json:"createdAt"`
		Metadata  map[string]string `json:"metadata,omitempty"`
	}

	for _, source := range sm.state.AllSources {
		serializableSources = append(serializableSources, struct {
			ID        string            `json:"id"`
			Number    int               `json:"number"`
			Type      string            `json:"type"`
			Path      string            `json:"path"`
			Title     string            `json:"title"`
			CreatedAt string            `json:"createdAt"`
			Metadata  map[string]string `json:"metadata,omitempty"`
		}{
			ID:        source.ID,
			Number:    source.Number,
			Type:      source.Type,
			Path:      source.Path,
			Title:     source.Title,
			CreatedAt: source.CreatedAt.Format(time.RFC3339),
			Metadata:  source.Metadata,
		})
	}

	tempState := struct {
		AllSources       interface{} `json:"allSources"`
		ActiveSourceIds  []string    `json:"activeSourceIds"`
		NextSourceNumber int         `json:"nextSourceNumber"`
	}{
		AllSources:       serializableSources,
		ActiveSourceIds:  sm.state.ActiveSourceIds,
		NextSourceNumber: sm.state.NextSourceNumber,
	}

	data, err := json.MarshalIndent(tempState, "", "  ")
	if err == nil {
		ioutil.WriteFile(SOURCES_FILE, data, 0644)
	}
}

func (sm *SourceManager) AddSource(sourceType string, path string) *KnowledgeSource {
	// Check if source already exists
	for _, source := range sm.state.AllSources {
		if source.Path == path {
			return &source
		}
	}

	// Find the lowest available number
	usedNumbers := make(map[int]bool)
	for _, source := range sm.state.AllSources {
		usedNumbers[source.Number] = true
	}

	availableNumber := 1
	for usedNumbers[availableNumber] {
		availableNumber++
	}

	id := generateUUID()

	newSource := KnowledgeSource{
		ID:        id,
		Number:    availableNumber,
		Type:      sourceType,
		Path:      path,
		Title:     filepath.Base(path),
		CreatedAt: time.Now(),
		Metadata:  make(map[string]string),
	}

	// Update nextSourceNumber to be the next available number
	if sm.state.NextSourceNumber <= availableNumber {
		sm.state.NextSourceNumber = availableNumber + 1
	}

	sm.state.AllSources = append(sm.state.AllSources, newSource)
	sm.saveSources()
	return &newSource
}

func (sm *SourceManager) ListSources() []KnowledgeSource {
	return sm.state.AllSources
}

func (sm *SourceManager) SetSourceActive(sourceId string, active bool) {
	if active {
		found := false
		for _, id := range sm.state.ActiveSourceIds {
			if id == sourceId {
				found = true
				break
			}
		}
		if !found {
			sm.state.ActiveSourceIds = append(sm.state.ActiveSourceIds, sourceId)
		}
	} else {
		var newActiveIds []string
		for _, id := range sm.state.ActiveSourceIds {
			if id != sourceId {
				newActiveIds = append(newActiveIds, id)
			}
		}
		sm.state.ActiveSourceIds = newActiveIds
	}
	sm.saveSources()
}

func (sm *SourceManager) SetActiveSources(sourceIds []string) {
	sm.state.ActiveSourceIds = sourceIds
	sm.saveSources()
}

func (sm *SourceManager) ClearActiveSources() {
	sm.state.ActiveSourceIds = []string{}
	sm.saveSources()
}

func (sm *SourceManager) GetActiveSources() []KnowledgeSource {
	var activeSources []KnowledgeSource
	for _, source := range sm.state.AllSources {
		for _, id := range sm.state.ActiveSourceIds {
			if source.ID == id {
				activeSources = append(activeSources, source)
				break
			}
		}
	}
	return activeSources
}

func (sm *SourceManager) GetActiveSourceIds() []string {
	return sm.state.ActiveSourceIds
}

func (sm *SourceManager) GetSourceByID(id string) *KnowledgeSource {
	// First try exact UUID match
	for _, source := range sm.state.AllSources {
		if source.ID == id {
			return &source
		}
	}

	// Then try partial UUID match (for backward compatibility with 4-char prefixes)
	for _, source := range sm.state.AllSources {
		if strings.HasPrefix(source.ID, id) {
			return &source
		}
	}

	// Then try sequential number match
	if number, err := strconv.Atoi(id); err == nil {
		for _, source := range sm.state.AllSources {
			if source.Number == number {
				return &source
			}
		}
	}

	return nil
}

func (sm *SourceManager) GetSourceByNumber(number int) *KnowledgeSource {
	for _, source := range sm.state.AllSources {
		if source.Number == number {
			return &source
		}
	}
	return nil
}

func (sm *SourceManager) GetSourceIdByNumber(number int) string {
	for _, source := range sm.state.AllSources {
		if source.Number == number {
			return source.ID
		}
	}
	return ""
}

func (sm *SourceManager) GetAllSources() []KnowledgeSource {
	return sm.state.AllSources
}

func (sm *SourceManager) UpdateSourceContent(sourceId string, newContent string) bool {
	source := sm.GetSourceByID(sourceId)
	if source == nil || source.Type != "file" {
		return false
	}

	// Write the new content back to the file
	err := ioutil.WriteFile(source.Path, []byte(newContent), 0644)
	if err != nil {
		return false
	}

	// After updating the file, we should re-ingest it to update the chunks
	return true
}

func (sm *SourceManager) GetSourceContent(sourceId string) string {
	source := sm.GetSourceByID(sourceId)
	if source == nil || source.Type != "file" {
		return ""
	}

	content, err := ioutil.ReadFile(source.Path)
	if err != nil {
		return ""
	}

	return string(content)
}

func (sm *SourceManager) RemoveSource(sourceId string) bool {
	// Try to find the source by UUID, partial UUID, or sequential number
	var matchingSource *KnowledgeSource

	// First try exact UUID match
	for _, source := range sm.state.AllSources {
		if source.ID == sourceId {
			matchingSource = &source
			break
		}
	}

	// Then try partial UUID match
	if matchingSource == nil {
		for _, source := range sm.state.AllSources {
			if strings.HasPrefix(source.ID, sourceId) {
				matchingSource = &source
				break
			}
		}
	}

	// Then try sequential number match
	if matchingSource == nil {
		if number, err := strconv.Atoi(sourceId); err == nil {
			for _, source := range sm.state.AllSources {
				if source.Number == number {
					matchingSource = &source
					break
				}
			}
		}
	}

	if matchingSource == nil {
		return false
	}

	initialLength := len(sm.state.AllSources)

	// Remove from all sources
	var newAllSources []KnowledgeSource
	for _, source := range sm.state.AllSources {
		if source.ID != matchingSource.ID {
			newAllSources = append(newAllSources, source)
		}
	}
	sm.state.AllSources = newAllSources

	// Also remove from active sources if it was active
	var newActiveIds []string
	for _, id := range sm.state.ActiveSourceIds {
		if id != matchingSource.ID {
			newActiveIds = append(newActiveIds, id)
		}
	}
	sm.state.ActiveSourceIds = newActiveIds

	if len(sm.state.AllSources) < initialLength {
		sm.saveSources()
		return true
	}
	return false
}

// Helper function to generate a UUID
func generateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant 10
	return hex.EncodeToString(bytes)
}
