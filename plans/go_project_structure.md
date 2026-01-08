# Go Implementation Plan for ЖЦП Parser

## Project Structure

```
zhcp-parser-go/
├── cmd/
│   └── zhcp-parser/
│       └── main.go
├── internal/
│   ├── parser/
│   │   ├── parser.go
│   │   └── parser_types.go
│   ├── parsers/
│   │   ├── pdf/
│   │   │   ├── pdf_parser.go
│   │   │   └── pdf_validator.go
│   │   ├── docx/
│   │   │   ├── docx_parser.go
│   │   │   └── docx_validator.go
│   │   └── text_preprocessor.go
│   ├── ai/
│   │   ├── llm_connector.go
│   │   ├── llm_providers/
│   │   │   ├── openai/
│   │   │   │   └── openai_provider.go
│   │   │   ├── anthropic/
│   │   │   │   └── anthropic_provider.go
│   │   │   ├── ollama/
│   │   │   │   └── ollama_provider.go
│   │   │   └── deepseek/
│   │   │       └── deepseek_provider.go
│   │   └── prompt_engineering.go
│   ├── transformers/
│   │   ├── json_transformer.go
│   │   └── data_enricher.go
│   ├── validators/
│   │   ├── data_validator.go
│   │   └── structure_validator.go
│   ├── config/
│   │   └── config.go
│   ├── errors/
│   │   └── errors.go
│   └── utils/
│       ├── logger.go
│       └── helpers.go
├── configs/
│   └── llm_config.yaml
├── docs/
│   └── api_spec.md
├── testdata/
│   ├── sample_project.pdf
│   └── sample_project.docx
├── go.mod
├── go.sum
├── README.md
└── USAGE_GUIDE.md
```

## Core Components Design

### 1. Main Parser (internal/parser/parser.go)

The main parser will orchestrate the entire pipeline:

```go
type ZhcpParser struct {
    config           *Config
    pdfExtractor     *PDFExtractor
    pdfValidator     *PDFValidator
    docxExtractor    *DOCXExtractor
    docxValidator    *DOCXValidator
    textPreprocessor *TextPreprocessor
    llmManager       *LLMManager
    promptManager    *PromptManager
    dataTransformer  *DataTransformer
    dataEnricher     *DataEnricher
    validationPipeline *ValidationPipeline
    errorHandler     *ErrorHandler
    logger           *Logger
}

func (p *ZhcpParser) ParseDocument(documentPath string, validate bool, enrich bool) (*ParseResult, error)
```

### 2. Document Parsers

#### PDF Parser (internal/parsers/pdf/pdf_parser.go)
- Uses libraries like `unidoc/unipdf` or `pdfcpu` for PDF parsing
- Handles text extraction while preserving structure
- Extracts tables and metadata
- Implements fallback mechanisms

#### DOCX Parser (internal/parsers/docx/docx_parser.go)
- Uses libraries like `baliance/sdcoffice` for DOCX parsing
- Extracts formatted text, tables, and metadata
- Preserves document structure

### 3. LLM Integration (internal/ai/llm_connector.go)

#### LLM Provider Interface
```go
type LLMProvider interface {
    Generate(ctx context.Context, prompt string, opts GenerationOptions) (*LLMResponse, error)
    GetCostEstimate(inputTokens, outputTokens int) float64
}

type LLMResponse struct {
    Content      string
    TokensUsed   TokenUsage
    Confidence   float64
    Model        string
    Timestamp    time.Time
    ParsedData   interface{}
}
```

#### Provider Implementations
- OpenAI Provider (using official Go SDK)
- Anthropic Provider (using official Go SDK)
- Ollama Provider (using HTTP client)
- DeepSeek Provider (using OpenAI-compatible interface)

### 4. Data Transformation (internal/transformers/json_transformer.go)

Handles conversion of LLM responses to structured project data:
- JSON parsing and validation
- Data normalization
- Date format standardization
- Confidence scoring

### 5. Data Validation (internal/validators/data_validator.go)

Multi-stage validation pipeline:
- Document content validation
- Structure validation
- Consistency validation
- Business rule validation

### 6. Error Handling (internal/errors/errors.go)

Comprehensive error system with:
- Error categories and severity levels
- Structured error information
- Error logging and tracking
- Suggested actions for each error type

### 7. Configuration (internal/config/config.go)

Handles configuration loading and management:
- YAML configuration files
- Environment variable substitution
- Provider configuration
- Fallback priorities
- Rate limiting settings

## Go-Specific Considerations

### Dependencies to be used:
- `golang.org/x/text` - for text processing and encoding
- `github.com/unidoc/unipdf/v3` - for PDF parsing
- `github.com/baliance/sdcoffice` - for DOCX parsing
- `github.com/sashabaranov/go-openai` - for OpenAI API
- `github.com/anthropics/anthropic-sdk-go` - for Anthropic API
- `gopkg.in/yaml.v3` - for YAML configuration
- `github.com/go-playground/validator/v10` - for data validation
- `github.com/rs/zerolog` - for structured logging
- `github.com/spf13/cobra` - for CLI commands

### Concurrency Considerations:
- Use goroutines for parallel document processing
- Implement connection pooling for API calls
- Use context for cancellation and timeouts
- Implement proper synchronization for shared resources

### Performance Optimizations:
- Implement caching for repeated operations
- Use efficient data structures
- Optimize memory usage for large documents
- Implement streaming for large file processing

## Data Models

### Project Structure (internal/parser/parser_types.go)
```go
type ProjectStructure struct {
    Project Project `json:"project"`
}

type Project struct {
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Phases      []Phase `json:"phases"`
    Metadata    Metadata `json:"metadata"`
}

type Phase struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    StartDate   string `json:"start_date,omitempty"`
    EndDate     string `json:"end_date,omitempty"`
    Tasks       []Task `json:"tasks"`
}

type Task struct {
    ID               string             `json:"id"`
    Name             string             `json:"name"`
    Description      string             `json:"description"`
    StartDate        string             `json:"start_date,omitempty"`
    EndDate          string             `json:"end_date,omitempty"`
    ResponsiblePersons []ResponsiblePerson `json:"responsible_persons"`
    Dependencies     []string           `json:"dependencies"`
    Status           string             `json:"status"`
}

type ResponsiblePerson struct {
    Name    string `json:"name"`
    Role    string `json:"role"`
    Contact string `json:"contact"`
}

type Metadata struct {
    SourceDocument    string            `json:"source_document"`
    ExtractionDate    string            `json:"extraction_date"`
    ConfidenceScore   float64           `json:"confidence_score"`
    ProcessingTime    float64           `json:"processing_time"`
    ValidationResults ValidationResults `json:"validation_results,omitempty"`
}
```

## CLI Interface

The command-line interface will support:
- Single document processing
- Batch processing
- Configuration options
- Output formatting options
- Verbose logging

```bash
# Basic usage
go run cmd/zhcp-parser/main.go parse <document_path>

# With options
go run cmd/zhcp-parser/main.go parse <document_path> --validate --enrich --output output.json

# Batch processing
go run cmd/zhcp-parser/main.go batch <directory_path>
```

## Testing Strategy

- Unit tests for each component
- Integration tests for the full pipeline
- Mock implementations for external services
- Sample documents for end-to-end testing
- Performance benchmarks