package llm

import (
	"fmt"
	"os"
)

func CreateLLMClient() LLMClient {
	useDeepSeek := os.Getenv("DEEPSEEK_API_KEY") != ""
	
	if useDeepSeek {
		fmt.Println("🚀 Using DeepSeek LLM")
		client, err := NewDeepSeekLLM()
		if err != nil {
			fmt.Printf("⚠️  Failed to initialize DeepSeek: %v\n", err)
			fmt.Println("🔄 Falling back to Mock LLM")
			return NewMockLLM()
		}
		return client
	} else {
		fmt.Println("🧪 Using Mock LLM (set DEEPSEEK_API_KEY to use DeepSeek)")
		return NewMockLLM()
	}
}