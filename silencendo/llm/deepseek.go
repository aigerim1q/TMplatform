package llm

import (
	stdcontext "context"
	"silencendo/context"
	"silencendo/ingestion"
	"os"
	
	"github.com/sashabaranov/go-openai"
)

type DeepSeekLLM struct {
	client *openai.Client
}

func NewDeepSeekLLM() (*DeepSeekLLM, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	model := os.Getenv("DEEPSEEK_MODEL")
	if model == "" {
		model = "deepseek-chat"
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	client := openai.NewClientWithConfig(config)
	
	return &DeepSeekLLM{
		client: client,
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
		Model:    os.Getenv("DEEPSEEK_MODEL"),
		Messages: openaiMessages,
	})
	if err != nil {
		return "", err
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
		Model: os.Getenv("DEEPSEEK_MODEL"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", err
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
		Model: os.Getenv("DEEPSEEK_MODEL"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}