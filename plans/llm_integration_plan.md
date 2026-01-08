# LLM Integration Modules Plan for Go Implementation

## Overview

This document outlines the plan for implementing LLM integration functionality in Go, supporting multiple providers (OpenAI, Anthropic, Ollama, DeepSeek) with fallback mechanisms.

## Architecture

### Core Interfaces

#### LLMProvider Interface
```go
type LLMProvider interface {
    Generate(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error)
    GetCostEstimate(inputTokens, outputTokens int) float64
    GetProviderType() ProviderType
}

type ProviderType string

const (
    OpenAIProvider    ProviderType = "openai"
    AnthropicProvider ProviderType = "anthropic"
    OllamaProvider    ProviderType = "ollama"
    DeepSeekProvider  ProviderType = "deepseek"
)

type GenerationOptions struct {
    Temperature float64
    MaxTokens   int
    Model       string
}

type LLMResponse struct {
    Content      string
    TokensUsed   TokenUsage
    Confidence   float64
    Model        string
    Timestamp    time.Time
    ParsedData   interface{} // Will be set after JSON parsing
}

type TokenUsage struct {
    Input  int
    Output int
    Total  int
}
```

### LLM Manager
```go
type LLMManager struct {
    config    *Config
    providers map[ProviderType]LLMProvider
    logger    *Logger
}

func (lm *LLMManager) GenerateWithFallback(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error)
func (lm *LLMManager) InitializeProviders() error
```

## Provider Implementations

### 1. OpenAI Provider (internal/ai/llm_providers/openai/openai_provider.go)

#### Dependencies
- `github.com/sashabaranov/go-openai` - OpenAI Go SDK

#### Implementation
```go
type OpenAIProvider struct {
    client   *openai.Client
    model    string
    logger   *Logger
    encoding *tokenizer.Tiktoken // for token counting
}

func NewOpenAIProvider(apiKey, model string) (*OpenAIProvider, error)
func (p *OpenAIProvider) Generate(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error)
func (p *OpenAIProvider) GetCostEstimate(inputTokens, outputTokens int) float64
func (p *OpenAIProvider) GetProviderType() ProviderType
func (p *OpenAIProvider) calculateConfidence(content string, usage TokenUsage) float64
```

#### Cost Calculation
- GPT-4 Turbo: $10/1M input tokens, $30/1M output tokens
- Implement based on OpenAI's current pricing

### 2. Anthropic Provider (internal/ai/llm_providers/anthropic/anthropic_provider.go)

#### Dependencies
- `github.com/anthropics/anthropic-sdk-go` - Anthropic Go SDK

#### Implementation
```go
type AnthropicProvider struct {
    client *anthropic.Client
    model  string
    logger *Logger
}

func NewAnthropicProvider(apiKey, model string) (*AnthropicProvider, error)
func (p *AnthropicProvider) Generate(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error)
func (p *AnthropicProvider) GetCostEstimate(inputTokens, outputTokens int) float64
func (p *AnthropicProvider) GetProviderType() ProviderType
func (p *AnthropicProvider) calculateConfidence(content string, usage TokenUsage) float64
```

#### Cost Calculation
- Claude 3 Sonnet: $3/1M input tokens, $15/1M output tokens
- Implement based on Anthropic's current pricing

### 3. Ollama Provider (internal/ai/llm_providers/ollama/ollama_provider.go)

#### Dependencies
- `github.com/ollama/ollama/api` - Ollama API client
- `golang.org/x/net/context` - for HTTP requests

#### Implementation
```go
type OllamaProvider struct {
    model    string
    baseURL  string
    client   *http.Client
    logger   *Logger
}

func NewOllamaProvider(model, baseURL string) (*OllamaProvider, error)
func (p *OllamaProvider) Generate(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error)
func (p *OllamaProvider) GetCostEstimate(inputTokens, outputTokens int) float64
func (p *OllamaProvider) GetProviderType() ProviderType
func (p *OllamaProvider) calculateConfidence(content string) float64
```

#### Notes
- Local models have zero cost
- Use HTTP API to communicate with Ollama server
- Implement proper error handling for local server unavailability

### 4. DeepSeek Provider (internal/ai/llm_providers/deepseek/deepseek_provider.go)

#### Dependencies
- `github.com/sashabaranov/go-openai` - DeepSeek is OpenAI-compatible

#### Implementation
```go
type DeepSeekProvider struct {
    client  *openai.Client
    model   string
    logger  *Logger
}

func NewDeepSeekProvider(apiKey, model string) (*DeepSeekProvider, error)
func (p *DeepSeekProvider) Generate(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error)
func (p *DeepSeekProvider) GetCostEstimate(inputTokens, outputTokens int) float64
func (p *DeepSeekProvider) GetProviderType() ProviderType
func (p *DeepSeekProvider) calculateConfidence(content string, usage TokenUsage) float64
```

#### Notes
- Uses OpenAI-compatible API
- Typically has free tier or minimal cost

## LLM Manager (internal/ai/llm_connector.go)

### Implementation
```go
type LLMManager struct {
    config           *Config
    providers        map[ProviderType]LLMProvider
    providerPriority []ProviderType
    logger           *Logger
}

func NewLLMManager(config *Config) (*LLMManager, error)
func (lm *LLMManager) InitializeProviders() error
func (lm *LLMManager) GenerateWithFallback(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error)
func (lm *LLMManager) initializeProvider(providerName string, providerConfig ProviderConfig) (LLMProvider, error)
```

### Fallback Strategy
1. Use provider priority from configuration
2. Try each provider in order until one succeeds
3. Log failures but continue to next provider
4. Return error only if all providers fail

### Configuration Integration
- Read provider settings from config file
- Handle environment variable substitution
- Validate provider configurations at startup

## Error Handling

### LLM-Specific Errors
```go
type LLMError struct {
    Provider ProviderType
    Message  string
    Details  map[string]interface{}
    Category ErrorCategory
}

func NewLLMError(provider ProviderType, message string, details map[string]interface{}) *LLMError
```

### Common Error Scenarios
- API key authentication failures
- Rate limiting
- Network connectivity issues
- Model unavailability
- Invalid responses from LLM

## Performance Considerations

### Rate Limiting
- Implement client-side rate limiting
- Respect provider-specific limits
- Use exponential backoff for retries
- Track token usage to stay within limits

### Connection Management
- Reuse HTTP connections
- Implement connection pooling
- Set appropriate timeouts
- Handle connection failures gracefully

### Caching
- Cache responses for identical prompts (if appropriate)
- Implement cache with TTL
- Consider cache invalidation strategies

## Security Considerations

### API Key Management
- Store API keys securely
- Use environment variables or secure vault
- Never log API keys
- Validate API keys before use

### Input Sanitization
- Sanitize prompts before sending to LLM
- Implement prompt injection protection
- Validate response content before processing

## Testing Strategy

### Unit Tests
- Test each provider implementation individually
- Test error handling scenarios
- Test cost calculation functions
- Test confidence scoring algorithms

### Integration Tests
- Test with actual API endpoints (using test keys)
- Test fallback mechanisms
- Test configuration loading
- Test concurrent usage scenarios

### Mock Implementations
- Create mock providers for testing
- Simulate different response scenarios
- Test error conditions and fallbacks

## Configuration Schema

### LLM Configuration (configs/llm_config.yaml)
```yaml
providers:
  openai:
    enabled: false
    api_key: "${OPENAI_API_KEY}"  # Use environment variable
    model: "gpt-4-turbo"
    temperature: 0.1
    max_tokens: 4096

  anthropic:
    enabled: false
    api_key: "${ANTHROPIC_API_KEY}"  # Use environment variable
    model: "claude-3-sonnet-20240229"
    temperature: 0.1
    max_tokens: 4096

  ollama:
    enabled: true  # Default to local model
    model: "llama3"
    base_url: "http://localhost:11434"
    temperature: 0.1
    max_tokens: 4096

  deepseek:
    enabled: false
    api_key: "${DEEPSEEK_API_KEY}"
    model: "deepseek-chat"
    temperature: 0.1
    max_tokens: 4096

provider_priority:
  - "deepseek"    # Primary provider
  - "ollama"      # Fallback 1
  - "openai"      # Fallback 2
  - "anthropic"   # Fallback 3

retry_settings:
  max_retries: 3
  backoff_factor: 1.0
  status_codes: [429, 502, 503, 504]

rate_limiting:
  requests_per_minute: 60
  tokens_per_minute: 100000
```

## Monitoring and Logging

### Metrics to Track
- API call success/failure rates
- Response times
- Token usage
- Cost estimation
- Provider usage distribution

### Logging Strategy
- Log provider selection decisions
- Log API call details (without sensitive data)
- Log performance metrics
- Log error details for debugging

This plan provides a comprehensive foundation for implementing LLM integration in Go with support for multiple providers and robust fallback mechanisms.