package llm

import (
	stdcontext "context"
	"silencendo/context"
	"silencendo/ingestion"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMClient interface {
	Generate(ctx stdcontext.Context, messages []Message) (string, error)
	Answer(ctx stdcontext.Context, question string, retrievedChunks []ingestion.Chunk, context context.ContextType) (string, error)
	Edit(ctx stdcontext.Context, inputText string, instruction string, retrievedChunks []ingestion.Chunk, context context.ContextType) (string, error)
}