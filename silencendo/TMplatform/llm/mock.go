package llm

import (
	stdcontext "context"
	botcontext "silencendo/context"
	"silencendo/ingestion"
)

type MockLLM struct{}

func NewMockLLM() *MockLLM {
	return &MockLLM{}
}

func (m *MockLLM) Generate(ctx stdcontext.Context, messages []Message) (string, error) {
	return "LLM response unavailable (mock mode). Enable DeepSeek for AI-generated answers.", nil
}

func (m *MockLLM) Answer(ctx stdcontext.Context, question string, retrievedChunks []ingestion.Chunk, context botcontext.ContextType) (string, error) {
	return "LLM response unavailable (mock mode). Enable DeepSeek for AI-generated answers.", nil
}

func (m *MockLLM) Edit(ctx stdcontext.Context, inputText string, instruction string, retrievedChunks []ingestion.Chunk, context botcontext.ContextType) (string, error) {
	return "LLM edit unavailable (mock mode). Enable DeepSeek for AI-generated edits.", nil
}
