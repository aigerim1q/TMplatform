package llm

import (
	"context"
	"silencendo/context"
	"silencendo/ingestion"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMClient interface {
	Generate(ctx context.Context, messages []Message) (string, error)
	Answer(ctx context.Context, question string, retrievedChunks []ingestion.Chunk, context context.ContextType) (string, error)
	Edit(ctx context.Context, inputText string, instruction string, retrievedChunks []ingestion.Chunk, context context.ContextType) (string, error)
}