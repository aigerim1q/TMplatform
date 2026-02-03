# Configuration Management Plan for Go Implementation

## Overview

This document outlines the plan for implementing configuration management functionality in Go, supporting YAML configuration files with environment variable substitution and provider-specific settings.

## Configuration Structure

### Core Configuration Types

#### Main Configuration
```go
type Config struct {
    Providers       map[string]ProviderConfig `yaml:"providers" json:"providers"`
    ProviderPriority []string                 `yaml:"provider_priority" json:"provider_priority"`
    RetrySettings   RetrySettings            `yaml:"retry_settings" json:"retry_settings"`
    RateLimiting    RateLimiting             `yaml:"rate_limiting" json:"rate_limiting"`
    ErrorHandling   ErrorHandlingConfig      `yaml:"error_handling" json:"error_handling"`
}

type ProviderConfig struct {
    Enabled     bool                   `yaml:"enabled" json:"enabled"`
    APIKey      string                 `yaml:"api_key" json:"api_key"`
    Model       string                 `yaml:"model" json:"model"`
    Temperature float64                `yaml:"temperature" json:"temperature"`
    MaxTokens   int                    `yaml:"max_tokens" json:"max_tokens"`
    BaseURL     string                 `yaml:"base_url,omitempty" json:"base_url,omitempty"`
    Details     map[string]interface{} `yaml:",inline" json:",omitempty"`
}
```

#### Retry Settings
```go
type RetrySettings struct {
    MaxRetries    int     `yaml:"max_retries" json:"max_retries"`
    BackoffFactor float64 `yaml:"backoff_factor" json:"backoff_factor"`
    StatusCodes   []int   `yaml:"status_codes" json:"status_codes"`
}
```

#### Rate Limiting
```go
type RateLimiting struct {
    RequestsPerMinute int `yaml:"requests_per_minute" json:"requests_per_minute"`
    TokensPerMinute   int `yaml:"tokens_per_minute" json:"tokens_per_minute"`
}
```

#### Error Handling Configuration
```go
type ErrorHandlingConfig struct {
    LogFile       string  `yaml:"log_file" json:"log_file"`
    MaxErrors     int     `yaml:"max_errors" json:"max_errors"`
    LogLevel      string  `yaml:"log_level" json:"log_level"`
    ErrorTolerance float64 `yaml:"error_tolerance" json:"error_tolerance"`
    RecoveryEnabled bool   `yaml:"recovery_enabled" json:"recovery_enabled"`
}
```

## Configuration Manager (internal/config/config.go)

### Core Structure
```go
type ConfigManager struct {
    configPath string
    config     *Config
    logger     *Logger
    mutex      sync.RWMutex  // For thread safety
}
```

### Constructor and Initialization
```go
func NewConfigManager(configPath string) (*ConfigManager, error)
func (cm *ConfigManager) LoadConfig() error
func (cm *ConfigManager) LoadDefaultConfig() *Config
func (cm *ConfigManager) SaveConfig(config *Config) error
```

### Configuration Loading
```go
func (cm *ConfigManager) LoadConfig() error {
    // Check if config file exists
    // If not, create default config
    // Load YAML content
    // Substitute environment variables
    // Validate configuration
    // Return error if validation fails
}
```

### Environment Variable Substitution
```go
func (cm *ConfigManager) substituteEnvVars(configContent []byte) []byte
func substituteEnvVarsRecursive(data interface{}) interface{}
```

## Configuration Validation

### Validation Logic
```go
func (cm *ConfigManager) ValidateConfig(config *Config) error
func (cm *ConfigManager) validateProviderConfig(providerName string, providerConfig ProviderConfig) error
func (cm *ConfigManager) validateProviderPriority(priority []string, availableProviders map[string]ProviderConfig) error
```

### Validation Rules
- Provider configurations must have required fields
- Provider priority list must reference existing providers
- Rate limiting values must be positive
- Retry settings must be within reasonable bounds
- API keys should follow expected format patterns

## Configuration Provider Interface

### Interface Definition
```go
type ConfigProvider interface {
    Load() (*Config, error)
    Save(config *Config) error
    Validate(config *Config) error
    WatchForChanges(callback func(*Config)) error
}
```

### YAML Configuration Provider
```go
type YAMLConfigProvider struct {
    filePath string
    logger   *Logger
}

func NewYAMLConfigProvider(filePath string) *YAMLConfigProvider
func (ycp *YAMLConfigProvider) Load() (*Config, error)
func (ycp *YAMLConfigProvider) Save(config *Config) error
func (ycp *YAMLConfigProvider) Validate(config *Config) error
```

## Environment Variable Handling

### Environment Variable Patterns
- `${VAR_NAME}` - Required environment variable
- `${VAR_NAME:default_value}` - Optional with default
- `${VAR_NAME?error_message}` - Required with error message

### Substitution Implementation
```go
type EnvVarSubstitutor struct {
    logger *Logger
}

func (evs *EnvVarSubstitutor) Substitute(content []byte) ([]byte, error)
func (evs *EnvVarSubstitutor) substituteString(input string) string
```

## Default Configuration

### Default Values
```go
func GetDefaultConfig() *Config {
    return &Config{
        Providers: map[string]ProviderConfig{
            "ollama": {
                Enabled:     true,
                Model:       "llama3",
                BaseURL:     "http://localhost:11434",
                Temperature: 0.1,
                MaxTokens:   4096,
            },
            "openai": {
                Enabled:     false,
                APIKey:      "${OPENAI_API_KEY}",
                Model:       "gpt-4-turbo",
                Temperature: 0.1,
                MaxTokens:   4096,
            },
            "anthropic": {
                Enabled:     false,
                APIKey:      "${ANTHROPIC_API_KEY}",
                Model:       "claude-3-sonnet-20240229",
                Temperature: 0.1,
                MaxTokens:   4096,
            },
            "deepseek": {
                Enabled:     false,
                APIKey:      "${DEEPSEEK_API_KEY}",
                Model:       "deepseek-chat",
                Temperature: 0.1,
                MaxTokens:   4096,
            },
        },
        ProviderPriority: []string{"ollama", "openai", "anthropic", "deepseek"},
        RetrySettings: RetrySettings{
            MaxRetries:    3,
            BackoffFactor: 1.0,
            StatusCodes:   []int{429, 502, 503, 504},
        },
        RateLimiting: RateLimiting{
            RequestsPerMinute: 60,
            TokensPerMinute:   100000,
        },
        ErrorHandling: ErrorHandlingConfig{
            LogFile:         "logs/errors.log",
            MaxErrors:       1000,
            LogLevel:        "INFO",
            ErrorTolerance:  0.1,
            RecoveryEnabled: true,
        },
    }
}
```

## Configuration Change Monitoring

### File Watcher
```go
type ConfigWatcher struct {
    filePath string
    callbacks []func(*Config)
    stopChan  chan struct{}
    logger    *Logger
}

func NewConfigWatcher(filePath string) *ConfigWatcher
func (cw *ConfigWatcher) Start() error
func (cw *ConfigWatcher) Stop()
func (cw *ConfigWatcher) AddCallback(callback func(*Config))
```

## Secure Configuration Handling

### API Key Security
- Never log API keys
- Validate API key format
- Support for encrypted configuration
- Secure temporary file handling

### Security Measures
```go
func (cm *ConfigManager) sanitizeConfigForLogging(config *Config) *Config
func validateAPIKeyFormat(providerName, apiKey string) bool
```

## Configuration Schema

### Expected YAML Structure (configs/llm_config.yaml)
```yaml
providers:
  ollama:
    enabled: true
    model: "llama3"
    base_url: "http://localhost:11434"
    temperature: 0.1
    max_tokens: 4096

  openai:
    enabled: false
    api_key: "${OPENAI_API_KEY}"
    model: "gpt-4-turbo"
    temperature: 0.1
    max_tokens: 4096

  anthropic:
    enabled: false
    api_key: "${ANTHROPIC_API_KEY}"
    model: "claude-3-sonnet-20240229"
    temperature: 0.1
    max_tokens: 4096

  deepseek:
    enabled: false
    api_key: "${DEEPSEEK_API_KEY}"
    model: "deepseek-chat"
    temperature: 0.1
    max_tokens: 4096

provider_priority:
  - "ollama"
  - "openai"
  - "anthropic"
  - "deepseek"

retry_settings:
  max_retries: 3
  backoff_factor: 1.0
  status_codes: [429, 502, 503, 504]

rate_limiting:
  requests_per_minute: 60
  tokens_per_minute: 100000

error_handling:
  log_file: "logs/errors.log"
  max_errors: 1000
  log_level: "INFO"
  error_tolerance: 0.1
  recovery_enabled: true
```

## Command-Line Configuration

### Configuration via CLI Flags
```go
type CLIConfig struct {
    ConfigFile   string
    LogLevel     string
    Provider     string
    Model        string
    Temperature  float64
    MaxTokens    int
}
```

### CLI Flag Handling
```go
func BindCLIConfigToViper(cliConfig CLIConfig, config *Config) error
```

## Testing Configuration Management

### Unit Tests
- Test configuration loading from YAML
- Test environment variable substitution
- Test configuration validation
- Test default configuration generation
- Test configuration saving

### Integration Tests
- Test configuration with actual environment variables
- Test configuration change detection
- Test configuration security measures
- Test concurrent configuration access

## Configuration Utilities

### Helper Functions
```go
func IsProviderEnabled(config *Config, providerName string) bool
func GetProviderConfig(config *Config, providerName string) (*ProviderConfig, bool)
func GetEnabledProviders(config *Config) []string
func ValidateProviderPriority(config *Config) error
```

## Error Handling for Configuration

### Configuration-Specific Errors
```go
type ConfigError struct {
    Type    string
    Message string
    Details map[string]interface{}
}

func NewConfigError(errorType, message string, details map[string]interface{}) *ConfigError
```

### Error Types
- ConfigLoadError: Error loading configuration file
- ConfigValidationError: Error validating configuration
- ConfigFormatError: Error in configuration format
- ConfigSecurityError: Security-related configuration error

This configuration management plan provides a robust foundation for handling application settings in the Go implementation, with support for environment variables, validation, and secure handling of sensitive information.