# Error Handling Modules Plan for Go Implementation

## Overview

This document outlines the plan for implementing comprehensive error handling functionality in Go, based on the existing Python implementation.

## Error Categories and Types

### Error Category Enum
```go
type ErrorCategory string

const (
    ErrorCategoryParsing           ErrorCategory = "parsing_error"
    ErrorCategoryValidation        ErrorCategory = "validation_error"
    ErrorCategoryLLM               ErrorCategory = "llm_error"
    ErrorCategoryTransformation    ErrorCategory = "transformation_error"
    ErrorCategoryConfiguration     ErrorCategory = "configuration_error"
    ErrorCategoryNetwork           ErrorCategory = "network_error"
    ErrorCategoryFile              ErrorCategory = "file_error"
    ErrorCategoryBusinessLogic     ErrorCategory = "business_logic_error"
    ErrorCategoryGeneral           ErrorCategory = "error"
)
```

### Error Severity Enum
```go
type ErrorSeverity string

const (
    ErrorSeverityInfo     ErrorSeverity = "info"
    ErrorSeverityWarning  ErrorSeverity = "warning"
    ErrorSeverityError    ErrorSeverity = "error"
    ErrorSeverityCritical ErrorSeverity = "critical"
)
```

## Base Error Types

### ErrorInfo Structure
```go
type ErrorInfo struct {
    ErrorID          string                 `json:"error_id"`
    Category         ErrorCategory          `json:"category"`
    Severity         ErrorSeverity          `json:"severity"`
    Message          string                 `json:"message"`
    Details          map[string]interface{} `json:"details"`
    Timestamp        time.Time              `json:"timestamp"`
    DocumentPath     *string                `json:"document_path,omitempty"`
    Component        *string                `json:"component,omitempty"`
    Traceback        *string                `json:"traceback,omitempty"`
    SuggestedAction  *string                `json:"suggested_action,omitempty"`
}
```

### Base Error Interface
```go
type ZhcpError interface {
    Error() string
    GetCategory() ErrorCategory
    GetDetails() map[string]interface{}
    GetTimestamp() time.Time
    GetErrorID() string
}
```

### Base Error Implementation
```go
type BaseError struct {
    Message   string
    Category  ErrorCategory
    Details   map[string]interface{}
    Timestamp time.Time
    ErrorID   string
}

func (e *BaseError) Error() string
func (e *BaseError) GetCategory() ErrorCategory
func (e *BaseError) GetDetails() map[string]interface{}
func (e *BaseError) GetTimestamp() time.Time
func (e *BaseError) GetErrorID() string
```

## Specific Error Types

### Parsing Error
```go
type ParsingError struct {
    *BaseError
    DocumentPath string
}

func NewParsingError(message, documentPath string, details map[string]interface{}) *ParsingError
```

### Validation Error
```go
type ValidationError struct {
    *BaseError
    Field string
    Value interface{}
}

func NewValidationError(message, field string, value interface{}) *ValidationError
```

### LLM Error
```go
type LLMError struct {
    *BaseError
    Provider string
}

func NewLLMError(message, provider string, details map[string]interface{}) *LLMError
```

### Transformation Error
```go
type TransformationError struct {
    *BaseError
}

func NewTransformationError(message string, details map[string]interface{}) *TransformationError
```

## Error Handler System (internal/errors/errors.go)

### ErrorHandler Structure
```go
type ErrorHandler struct {
    logFile    string
    maxErrors  int
    errorLog   []*ErrorInfo
    logger     *Logger
    mutex      sync.RWMutex  // For thread safety
}

func NewErrorHandler(logFile string, maxErrors int) *ErrorHandler
```

### Core Methods
```go
func (eh *ErrorHandler) HandleError(err error, documentPath, component string) *ErrorInfo
func (eh *ErrorHandler) createErrorInfo(err error, documentPath, component string) *ErrorInfo
func (eh *ErrorHandler) logError(errorInfo *ErrorInfo) error
func (eh *ErrorHandler) GetErrorSummary() map[string]interface{}
func (eh *ErrorHandler) determineSeverity(category ErrorCategory) ErrorSeverity
func (eh *ErrorHandler) getSuggestedAction(category ErrorCategory) string
```

### Error Handling Logic
```go
// handleError processes an error and returns comprehensive error information
func (eh *ErrorHandler) HandleError(err error, documentPath, component string) *ErrorInfo

// createErrorInfo creates error information from any error type
func (eh *ErrorHandler) createErrorInfo(err error, documentPath, component string) *ErrorInfo

// logError writes error information to the log file
func (eh *ErrorHandler) logError(errorInfo *ErrorInfo) error

// GetErrorSummary returns a summary of recent errors
func (eh *ErrorHandler) GetErrorSummary() map[string]interface{}
```

## Error Logging

### Log Format
```go
type LogEntry struct {
    ErrorID      string                 `json:"error_id"`
    Category     ErrorCategory          `json:"category"`
    Severity     ErrorSeverity          `json:"severity"`
    Message      string                 `json:"message"`
    Timestamp    string                 `json:"timestamp"`
    DocumentPath *string                `json:"document_path,omitempty"`
    Component    *string                `json:"component,omitempty"`
    Details      map[string]interface{} `json:"details"`
}
```

### File-based Logging
- Store errors in JSON format
- Implement log rotation
- Handle concurrent writes safely
- Include timestamp and error context

## Error Recovery and Fallbacks

### Recovery Strategies
- Graceful degradation when errors occur
- Fallback to alternative providers/strategies
- Partial result generation when possible
- Continue processing when non-critical errors occur

### Fallback Implementation
```go
type FallbackHandler struct {
    errorHandler *ErrorHandler
    logger       *Logger
}

func (fh *FallbackHandler) ExecuteWithFallback(
    operations []func() error,
    fallbackMessage string,
) error
```

## Integration with Other Modules

### Document Parsing Error Integration
```go
// In PDF parser
func (p *PDFExtractor) ExtractText(pdfPath string) (*PDFExtractionResult, error) {
    // ... implementation
    if err != nil {
        errorInfo := p.errorHandler.HandleError(err, pdfPath, "PDFExtractor")
        return nil, NewParsingError(errorInfo.Message, pdfPath, errorInfo.Details)
    }
}
```

### LLM Error Integration
```go
// In LLM provider
func (p *OpenAIProvider) Generate(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error) {
    // ... implementation
    if err != nil {
        errorInfo := p.errorHandler.HandleError(err, "", "OpenAIProvider")
        return nil, NewLLMError(errorInfo.Message, "openai", errorInfo.Details)
    }
}
```

### Validation Error Integration
```go
// In validation
func (v *StructureValidator) ValidateStructure(structure *ProjectStructure) *ValidationResult {
    // ... implementation
    if !isValid {
        errorInfo := v.errorHandler.HandleError(
            NewValidationError("Structure validation failed", "structure", structure), 
            "", "StructureValidator")
        // Add to validation result
    }
}
```

## Error Context Propagation

### Context with Error Handling
```go
type ErrorContext struct {
    ctx    context.Context
    errors []error
    mu     sync.Mutex
}

func NewErrorContext(ctx context.Context) *ErrorContext
func (ec *ErrorContext) AddError(err error)
func (ec *ErrorContext) GetErrors() []error
func (ec *ErrorContext) HasErrors() bool
```

## Error Reporting

### Error Summary Structure
```go
type ErrorSummary struct {
    TotalErrors   int                           `json:"total_errors"`
    ByCategory    map[string]int               `json:"by_category"`
    BySeverity    map[string]int               `json:"by_severity"`
    RecentErrors  []ErrorInfo                  `json:"recent_errors"`
    ErrorRate     float64                      `json:"error_rate"`
    LastError     *ErrorInfo                   `json:"last_error,omitempty"`
}
```

## Error Metrics and Monitoring

### Metrics Collection
- Count errors by category and severity
- Track error rates over time
- Monitor specific error patterns
- Alert on critical errors

### Monitoring Integration
```go
type ErrorMetrics struct {
    totalErrors    int64
    categoryCounts map[ErrorCategory]int64
    severityCounts map[ErrorSeverity]int64
    mu             sync.RWMutex
}

func (em *ErrorMetrics) RecordError(err *ErrorInfo)
func (em *ErrorMetrics) GetMetrics() map[string]interface{}
```

## Thread Safety

### Concurrency Considerations
- Use mutexes for shared error log access
- Ensure thread-safe error handling
- Handle concurrent error logging
- Prevent race conditions in error counting

## Error Recovery Patterns

### Retry Logic
```go
type RetryConfig struct {
    MaxRetries    int
    BackoffFactor float64
    StatusCodes   []int
}

func WithRetry(ctx context.Context, fn func() error, config RetryConfig) error
```

### Circuit Breaker Pattern
```go
type CircuitBreaker struct {
    maxFailures int
    resetTime   time.Duration
    state       int32  // 0: closed, 1: open, 2: half-open
    failureCount int32
    lastFailure time.Time
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Call(ctx context.Context, operation func() error) error
```

## Testing Error Handling

### Unit Tests for Error Handling
- Test error creation and categorization
- Test error logging functionality
- Test error recovery mechanisms
- Test thread safety of error handling

### Error Scenarios to Test
- All types of errors (parsing, validation, LLM, etc.)
- Error with and without context information
- Concurrent error handling
- Error log file handling
- Error recovery and fallback mechanisms

## Configuration for Error Handling

### Error Handling Configuration
```go
type ErrorHandlingConfig struct {
    LogFile       string `yaml:"log_file"`
    MaxErrors     int    `yaml:"max_errors"`
    LogLevel      string `yaml:"log_level"`
    ErrorTolerance float64 `yaml:"error_tolerance"`
    RecoveryEnabled bool  `yaml:"recovery_enabled"`
}
```

This comprehensive error handling plan ensures robust error management throughout the application, with proper categorization, logging, recovery mechanisms, and monitoring capabilities.