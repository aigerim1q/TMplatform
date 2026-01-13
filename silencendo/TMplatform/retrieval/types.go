package retrieval

import "silencendo/ingestion"

type RetrievalResult struct {
	Chunks []ingestion.Chunk `json:"chunks"`
	Scores []float64         `json:"scores"`
}

type RetrievalOptions struct {
	TopK     int     `json:"topK"`
	MinScore float64 `json:"minScore"`
}
