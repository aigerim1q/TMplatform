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

	fmt.Println("🧪 Using Mock LLM (set DEEPSEEK_API_KEY to enable DeepSeek)")
	return NewMockLLM()
}
