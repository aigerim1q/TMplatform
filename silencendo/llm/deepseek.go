package llm

import (
	stdcontext "context"
	"fmt"
	"os"
	"strings"

	"silencendo/context"
	"silencendo/ingestion"

	"github.com/sashabaranov/go-openai"
)

type DeepSeekLLM struct {
	client *openai.Client
	model  string
}

const groundedAnswerSystem = "You are a retrieval QA agent. Use ONLY the provided context. If the context is empty or missing the answer, reply exactly: 'I don't have enough information to answer that yet. Please add sources or ingest them.' Do not invent names, people, or facts."

const groundedEditSystem = "You are an editing assistant. Use ONLY the provided context plus the input text. Do not invent names, people, or facts. If context is empty, operate only on the input text and instruction without adding new factual details."

func NewDeepSeekLLM() (*DeepSeekLLM, error) {
	apiKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	if apiKey == "" {
		return nil, fmt.Errorf("DEEPSEEK_API_KEY is required")
	}

	baseURL := strings.TrimSpace(os.Getenv("DEEPSEEK_BASE_URL"))
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}

	model := strings.TrimSpace(os.Getenv("DEEPSEEK_MODEL"))
	if model == "" {
		model = "deepseek-chat"
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	client := openai.NewClientWithConfig(config)

	return &DeepSeekLLM{
		client: client,
		model:  model,
	}, nil
}

func (d *DeepSeekLLM) Generate(ctx stdcontext.Context, messages []Message) (string, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	resp, err := d.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    d.model,
		Messages: openaiMessages,
	})
	if err != nil {
		// Try fallback to OpenAI if DeepSeek fails
		fmt.Printf("⚠️  DeepSeek failed: %v. Trying OpenAI fallback...\n", err)
		openaiClient, err := NewOpenAILLM()
		if err != nil {
			return "", fmt.Errorf("deepseek failed (%v) and OpenAI fallback failed (%v)", err, err)
		}
		return openaiClient.Generate(ctx, messages)
	}

	if len(resp.Choices) == 0 {
		// Try fallback to OpenAI if DeepSeek returns no choices
		fmt.Printf("⚠️  DeepSeek returned no choices. Trying OpenAI fallback...\n")
		openaiClient, err := NewOpenAILLM()
		if err != nil {
			return "", fmt.Errorf("deepseek returned no choices and OpenAI fallback failed (%v)", err)
		}
		return openaiClient.Generate(ctx, messages)
	}

	return resp.Choices[0].Message.Content, nil
}

func (d *DeepSeekLLM) Answer(ctx stdcontext.Context, question string, retrievedChunks []ingestion.Chunk, context context.ContextType) (string, error) {
	var prompt string

	if len(retrievedChunks) > 0 {
		// Build context from retrieved chunks
		var contextText string
		for _, chunk := range retrievedChunks {
			contextText += chunk.Text + "\n\n"
		}

		prompt = "Based on the following information, please answer the question:\n\n" +
			"Context:\n" + contextText + "\n\n" +
			"Question: " + question + "\n\n" +
			"Answer:"
	} else {
		prompt = "Please answer the following question:\n\n" +
			"Question: " + question + "\n\n" +
			"Answer:"
	}

	resp, err := d.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: d.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: groundedAnswerSystem},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		// Try fallback to OpenAI if DeepSeek fails
		fmt.Printf("⚠️  DeepSeek failed during answer: %v. Trying OpenAI fallback...\n", err)
		openaiClient, err := NewOpenAILLM()
		if err != nil {
			return "", fmt.Errorf("deepseek failed (%v) and OpenAI fallback failed (%v)", err, err)
		}
		return openaiClient.Answer(ctx, question, retrievedChunks, context)
	}

	if len(resp.Choices) == 0 {
		// Try fallback to OpenAI if DeepSeek returns no choices
		fmt.Printf("⚠️  DeepSeek returned no choices during answer. Trying OpenAI fallback...\n")
		openaiClient, err := NewOpenAILLM()
		if err != nil {
			return "", fmt.Errorf("deepseek returned no choices and OpenAI fallback failed (%v)", err)
		}
		return openaiClient.Answer(ctx, question, retrievedChunks, context)
	}

	return resp.Choices[0].Message.Content, nil
}

func (d *DeepSeekLLM) Edit(ctx stdcontext.Context, inputText string, instruction string, retrievedChunks []ingestion.Chunk, context context.ContextType) (string, error) {
	var prompt string

	if len(retrievedChunks) > 0 {
		// Build context from retrieved chunks
		var contextText string
		for _, chunk := range retrievedChunks {
			contextText += chunk.Text + "\n\n"
		}

		prompt = "Using the provided context and following the instruction, please edit the input text:\n\n" +
			"Context:\n" + contextText + "\n\n" +
			"Instruction: " + instruction + "\n\n" +
			"Input Text: " + inputText + "\n\n" +
			"Edited Text:"
	} else {
		prompt = "Please edit the following text according to the instruction:\n\n" +
			"Instruction: " + instruction + "\n\n" +
			"Input Text: " + inputText + "\n\n" +
			"Edited Text:"
	}

	resp, err := d.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: d.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: groundedEditSystem},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		// Try fallback to OpenAI if DeepSeek fails
		fmt.Printf("⚠️  DeepSeek failed during edit: %v. Trying OpenAI fallback...\n", err)
		openaiClient, err := NewOpenAILLM()
		if err != nil {
			return "", fmt.Errorf("deepseek failed (%v) and OpenAI fallback failed (%v)", err, err)
		}
		return openaiClient.Edit(ctx, inputText, instruction, retrievedChunks, context)
	}

	if len(resp.Choices) == 0 {
		// Try fallback to OpenAI if DeepSeek returns no choices
		fmt.Printf("⚠️  DeepSeek returned no choices during edit. Trying OpenAI fallback...\n")
		openaiClient, err := NewOpenAILLM()
		if err != nil {
			return "", fmt.Errorf("deepseek returned no choices and OpenAI fallback failed (%v)", err)
		}
		return openaiClient.Edit(ctx, inputText, instruction, retrievedChunks, context)
	}

	return resp.Choices[0].Message.Content, nil
}
