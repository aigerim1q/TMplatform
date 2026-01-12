package llm

import "testing"

func TestNewDeepSeekLLMRequiresAPIKey(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "")
	t.Setenv("DEEPSEEK_MODEL", "")
	t.Setenv("DEEPSEEK_BASE_URL", "")

	if _, err := NewDeepSeekLLM(); err == nil {
		t.Fatalf("expected error when DEEPSEEK_API_KEY is missing")
	}
}

func TestNewDeepSeekLLMUsesDefaults(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "test-key")
	t.Setenv("DEEPSEEK_BASE_URL", "http://example.com")
	t.Setenv("DEEPSEEK_MODEL", "")

	client, err := NewDeepSeekLLM()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.model != "deepseek-chat" {
		t.Fatalf("expected default model 'deepseek-chat', got %q", client.model)
	}
}
