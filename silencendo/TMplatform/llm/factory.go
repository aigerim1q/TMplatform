package llm

import (
	"fmt"
	"os"
	"strings"
)

func CreateLLMClient() LLMClient {
	apiKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	if apiKey != "" {
		client, err := NewDeepSeekLLM()
		if err == nil {
			fmt.Println("🚀 Using DeepSeek LLM")
			return client
		}
		fmt.Printf("⚠️  Failed to initialize DeepSeek: %v\n", err)
	}

	// Try OpenAI as fallback
	openaiAPIKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if openaiAPIKey != "" {
		client, err := NewOpenAILLM()
		if err == nil {
			fmt.Println("🚀 Using OpenAI LLM as fallback")
			return client
		}
		fmt.Printf("⚠️  Failed to initialize OpenAI: %v\n", err)
	}

	fmt.Println("🧪 Using Mock LLM (set DEEPSEEK_API_KEY or OPENAI_API_KEY to enable real LLM)")
	return NewMockLLM()
}
