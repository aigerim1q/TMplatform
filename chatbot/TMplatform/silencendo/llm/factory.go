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
		} else {
			return client
		}
	}

	// Try OpenAI as fallback
	useOpenAI := os.Getenv("OPENAI_API_KEY") != ""
	if useOpenAI {
		fmt.Println("🚀 Using OpenAI LLM as fallback")
		client, err := NewOpenAILLM()
		if err != nil {
			fmt.Printf("⚠️  Failed to initialize OpenAI: %v\n", err)
		} else {
			return client
		}
	}

	fmt.Println("🧪 Using Mock LLM (set DEEPSEEK_API_KEY or OPENAI_API_KEY to use real LLM)")
	return NewMockLLM()
}
