package llm

import (
	stdcontext "context"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"

	botcontext "silencendo/context"
	"silencendo/ingestion"
)

type DeepSeekLLM struct {
	client *openai.Client
}

const groundedAnswerSystem = "You are a retrieval QA agent. Use ONLY the provided context. If the context is empty or does not contain the answer, reply exactly: 'I don't have enough information to answer that yet. Please add sources or ingest them.' Do not invent names, people, or facts."

const groundedEditSystem = "You are an editing assistant. Use ONLY the provided context plus the given input text. Do not invent names, people, or facts. If context is empty, operate only on the input text and instruction without adding new factual details."

func NewDeepSeekLLM() (*DeepSeekLLM, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://api.deepseek.com"
	}

	model := os.Getenv("DEEPSEEK_MODEL")
	if strings.TrimSpace(model) == "" {
		model = "deepseek-chat"
	}

	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = baseURL

	client := openai.NewClientWithConfig(cfg)
	return &DeepSeekLLM{client: client}, nil
}

func (d *DeepSeekLLM) Generate(ctx stdcontext.Context, messages []Message) (string, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content}
	}

	resp, err := d.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    os.Getenv("DEEPSEEK_MODEL"),
		Messages: openaiMessages,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}

func (d *DeepSeekLLM) Answer(ctx stdcontext.Context, question string, retrievedChunks []ingestion.Chunk, context botcontext.ContextType) (string, error) {
	prompt := buildAnswerPrompt(question, retrievedChunks)
	resp, err := d.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: os.Getenv("DEEPSEEK_MODEL"),
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: groundedAnswerSystem},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}

func (d *DeepSeekLLM) Edit(ctx stdcontext.Context, inputText, instruction string, retrievedChunks []ingestion.Chunk, context botcontext.ContextType) (string, error) {
	prompt := buildEditPrompt(inputText, instruction, retrievedChunks)
	resp, err := d.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: os.Getenv("DEEPSEEK_MODEL"),
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: groundedEditSystem},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}

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
