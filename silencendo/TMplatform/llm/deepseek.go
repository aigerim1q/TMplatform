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
		Model:    os.Getenv("DEEPSEEK_MODEL"),
		Messages: []openai.ChatCompletionMessage{{Role: "user", Content: prompt}},
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
		Model:    os.Getenv("DEEPSEEK_MODEL"),
		Messages: []openai.ChatCompletionMessage{{Role: "user", Content: prompt}},
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
	if len(retrievedChunks) == 0 {
		return "Please answer the following question:\n\nQuestion: " + question + "\n\nAnswer:"
	}

	var b strings.Builder
	b.WriteString("Based on the following information, please answer the question:\n\nContext:\n")
	for _, chunk := range retrievedChunks {
		b.WriteString(chunk.Text)
		b.WriteString("\n\n")
	}
	b.WriteString("Question: ")
	b.WriteString(question)
	b.WriteString("\n\nAnswer:")
	return b.String()
}

func buildEditPrompt(inputText, instruction string, retrievedChunks []ingestion.Chunk) string {
	if len(retrievedChunks) == 0 {
		return "Please edit the following text according to the instruction:\n\nInstruction: " + instruction + "\n\nInput Text: " + inputText + "\n\nEdited Text:"
	}

	var b strings.Builder
	b.WriteString("Using the provided context and following the instruction, please edit the input text:\n\nContext:\n")
	for _, chunk := range retrievedChunks {
		b.WriteString(chunk.Text)
		b.WriteString("\n\n")
	}
	b.WriteString("Instruction: ")
	b.WriteString(instruction)
	b.WriteString("\n\nInput Text: ")
	b.WriteString(inputText)
	b.WriteString("\n\nEdited Text:")
	return b.String()
}
