package retrieval

import (
	"math"
	"regexp"
	"silencendo/ingestion"
	"strings"
)

type Retriever struct {
	chunkManager *ingestion.ChunkManager
	options      RetrievalOptions
}

func NewRetriever(chunkManager *ingestion.ChunkManager, options *RetrievalOptions) *Retriever {
	opts := RetrievalOptions{
		TopK:     5,
		MinScore: 0.01,
	}

	if options != nil {
		if options.TopK > 0 {
			opts.TopK = options.TopK
		}
		if options.MinScore >= 0 {
			opts.MinScore = options.MinScore
		}
	}

	return &Retriever{
		chunkManager: chunkManager,
		options:      opts,
	}
}

func (r *Retriever) Retrieve(query string, sourceIds []string) RetrievalResult {
	// Get all chunks or filter by source IDs
	var chunks []ingestion.Chunk
	if len(sourceIds) > 0 {
		for _, sourceId := range sourceIds {
			chunks = append(chunks, r.chunkManager.GetChunksBySource(sourceId)...)
		}
	} else {
		chunks = r.chunkManager.GetAllChunks()
	}

	if len(chunks) == 0 {
		return RetrievalResult{Chunks: []ingestion.Chunk{}, Scores: []float64{}}
	}

	// Calculate similarity scores for the query against all chunks
	scores := make([]float64, len(chunks))
	for i, chunk := range chunks {
		scores[i] = r.calculateSimilarity(query, chunk.Text)
	}

	// Create pairs of chunks and scores, then sort by score in descending order
	type chunkScorePair struct {
		chunk ingestion.Chunk
		score float64
	}

	chunkScorePairs := make([]chunkScorePair, len(chunks))
	for i := range chunks {
		chunkScorePairs[i] = chunkScorePair{chunk: chunks[i], score: scores[i]}
	}

	// Filter by minimum score, but if no chunks pass the filter and we have chunks,
	// return at least the highest scoring chunk to avoid empty results
	var filteredPairs []chunkScorePair
	for _, pair := range chunkScorePairs {
		if pair.score >= r.options.MinScore {
			filteredPairs = append(filteredPairs, pair)
		}
	}

	if len(filteredPairs) == 0 && len(chunkScorePairs) > 0 {
		// If no chunks meet the threshold, return the best scoring chunk anyway
		// This prevents empty results when the query is not similar to any chunk
		bestPair := chunkScorePairs[0]
		for _, pair := range chunkScorePairs[1:] {
			if pair.score > bestPair.score {
				bestPair = pair
			}
		}
		chunkScorePairs = []chunkScorePair{bestPair}
	} else {
		chunkScorePairs = filteredPairs
	}

	// Sort by score and take topK
	for i := 0; i < len(chunkScorePairs)-1; i++ {
		for j := i + 1; j < len(chunkScorePairs); j++ {
			if chunkScorePairs[i].score < chunkScorePairs[j].score {
				chunkScorePairs[i], chunkScorePairs[j] = chunkScorePairs[j], chunkScorePairs[i]
			}
		}
	}

	if len(chunkScorePairs) > r.options.TopK {
		chunkScorePairs = chunkScorePairs[:r.options.TopK]
	}

	// Extract chunks and scores
	resultChunks := make([]ingestion.Chunk, len(chunkScorePairs))
	resultScores := make([]float64, len(chunkScorePairs))

	for i, pair := range chunkScorePairs {
		resultChunks[i] = pair.chunk
		resultScores[i] = pair.score
	}

	return RetrievalResult{
		Chunks: resultChunks,
		Scores: resultScores,
	}
}

func (r *Retriever) calculateSimilarity(query string, text string) float64 {
	// Enhanced similarity calculation with multiple strategies
	queryWords := r.tokenize(strings.ToLower(query))
	textWords := r.tokenize(strings.ToLower(text))

	if len(queryWords) == 0 || len(textWords) == 0 {
		return 0
	}

	// Calculate Jaccard similarity (intersection over union)
	queryWordSet := make(map[string]bool)
	textWordSet := make(map[string]bool)

	for _, word := range queryWords {
		queryWordSet[word] = true
	}
	for _, word := range textWords {
		textWordSet[word] = true
	}

	overlap := 0
	for word := range queryWordSet {
		if textWordSet[word] {
			overlap++
		}
	}

	jaccard := float64(overlap) / math.Max(float64(len(queryWordSet)), float64(len(textWordSet)))

	// Also calculate overlap ratio relative to query size (for queries that are shorter than text)
	queryOverlapRatio := float64(overlap) / math.Max(float64(len(queryWordSet)), 1)

	// Use the higher of the two scores to handle both long text vs short query and vice versa
	return math.Max(jaccard, queryOverlapRatio)
}

func (r *Retriever) tokenize(text string) []string {
	// Simple tokenization: split on non-word characters and filter empty strings
	re := regexp.MustCompile(`\W+`)
	tokens := re.Split(text, -1)

	var result []string
	for _, token := range tokens {
		if token != "" {
			result = append(result, token)
		}
	}
	return result
}
