package llm

import (
	stdcontext "context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"

	botcontext "silencendo/context"
	"silencendo/ingestion"
)

type OpenAILLM struct {
	client *openai.Client
	model  string
}

func NewOpenAILLM() (*OpenAILLM, error) {
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	client := openai.NewClient(apiKey)

	return &OpenAILLM{
		client: client,
		model:  model,
	}, nil
}

func (o *OpenAILLM) Generate(ctx stdcontext.Context, messages []Message) (string, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: openaiMessages,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAILLM) Answer(ctx stdcontext.Context, question string, retrievedChunks []ingestion.Chunk, context botcontext.ContextType) (string, error) {
	prompt := buildAnswerPrompt(question, retrievedChunks)

	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: o.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: groundedAnswerSystem},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAILLM) Edit(ctx stdcontext.Context, inputText, instruction string, retrievedChunks []ingestion.Chunk, context botcontext.ContextType) (string, error) {
	prompt := buildEditPrompt(inputText, instruction, retrievedChunks)

	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: o.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: groundedEditSystem},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return resp.Choices[0].Message.Content, nil
}
