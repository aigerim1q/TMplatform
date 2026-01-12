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

// These helper functions are copied from deepseek.go to maintain consistency
func buildAnswerPrompt(question string, retrievedChunks []ingestion.Chunk) string {
	var b strings.Builder
	b.WriteString("You must answer using only the provided context. If the context lacks the answer, say you do not have enough information.\n\nContext:\n")
	if len(retrievedChunks) == 0 {
		b.WriteString("<no context provided>\n\n")
	} else {
		for _, chunk := range retrievedChunks {
			b.WriteString(chunk.Text)
			b.WriteString("\n\n")
		}
	}
	b.WriteString("Question: ")
	b.WriteString(question)
	b.WriteString("\n\nAnswer:")
	return b.String()
}

func buildEditPrompt(inputText, instruction string, retrievedChunks []ingestion.Chunk) string {
	var b strings.Builder
	b.WriteString("Edit only using the provided context plus the input text. Do not invent facts.\n\nContext:\n")
	if len(retrievedChunks) == 0 {
		b.WriteString("<no context provided>\n\n")
	} else {
		for _, chunk := range retrievedChunks {
			b.WriteString(chunk.Text)
			b.WriteString("\n\n")
		}
	}
	b.WriteString("Instruction: ")
	b.WriteString(instruction)
	b.WriteString("\n\nInput Text: ")
	b.WriteString(inputText)
	b.WriteString("\n\nEdited Text:")
	return b.String()
}
