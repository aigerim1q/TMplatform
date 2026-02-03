package ingestion

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const CHUNKS_FILE = ".bot/chunks.json"

type ChunkStorage struct {
	Chunks []Chunk `json:"chunks"`
}

type ChunkManager struct {
	storage ChunkStorage
}

func NewChunkManager() *ChunkManager {
	cm := &ChunkManager{
		storage: ChunkStorage{
			Chunks: []Chunk{},
		},
	}
	cm.loadChunks()
	return cm
}

func (cm *ChunkManager) ensureBotDir() {
	if _, err := os.Stat(".bot"); os.IsNotExist(err) {
		os.MkdirAll(".bot", 0755)
	}
}

func (cm *ChunkManager) loadChunks() {
	cm.ensureBotDir()

	if _, err := os.Stat(CHUNKS_FILE); err == nil {
		data, err := ioutil.ReadFile(CHUNKS_FILE)
		if err == nil {
			err = json.Unmarshal(data, &cm.storage)
			if err != nil {
				// Use default state if loading fails
				cm.saveChunks()
			}
		} else {
			// Use default state if loading fails
			cm.saveChunks()
		}
	} else {
		// Initialize with empty state
		cm.saveChunks()
	}
}

func (cm *ChunkManager) saveChunks() {
	cm.ensureBotDir()
	data, err := json.MarshalIndent(cm.storage, "", "  ")
	if err == nil {
		ioutil.WriteFile(CHUNKS_FILE, data, 0644)
	}
}

func (cm *ChunkManager) AddChunks(chunks []Chunk) {
	cm.storage.Chunks = append(cm.storage.Chunks, chunks...)
	cm.saveChunks()
}

func (cm *ChunkManager) GetChunksBySource(sourceId string) []Chunk {
	var result []Chunk
	for _, chunk := range cm.storage.Chunks {
		if chunk.SourceID == sourceId {
			result = append(result, chunk)
		}
	}
	return result
}

func (cm *ChunkManager) GetAllChunks() []Chunk {
	return cm.storage.Chunks
}

func (cm *ChunkManager) RemoveChunksBySource(sourceId string) {
	var newChunks []Chunk
	for _, chunk := range cm.storage.Chunks {
		if chunk.SourceID != sourceId {
			newChunks = append(newChunks, chunk)
		}
	}
	cm.storage.Chunks = newChunks
	cm.saveChunks()
}

func (cm *ChunkManager) ClearChunks() {
	cm.storage.Chunks = []Chunk{}
	cm.saveChunks()
}

func (cm *ChunkManager) GetChunksByIds(ids []string) []Chunk {
	var result []Chunk
	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	for _, chunk := range cm.storage.Chunks {
		if idMap[chunk.ChunkID] {
			result = append(result, chunk)
		}
	}
	return result
}
