package ingestion

type Chunk struct {
	ChunkID  string `json:"chunkId"`
	SourceID string `json:"sourceId"`
	Text     string `json:"text"`
	Metadata struct {
		Position int `json:"position"`
		Length   int `json:"length"`
	} `json:"metadata"`
}

type ChunkingOptions struct {
	ChunkSize int `json:"chunkSize"`
	Overlap   int `json:"overlap"`
}

type IngestionResult struct {
	Chunks  []Chunk `json:"chunks"`
	SourceID string `json:"sourceId"`
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
}