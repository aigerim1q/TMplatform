# ЖЦП Parser - AI-Powered Project Lifecycle Document Parser

An AI-powered module that automatically extracts project structure information from PDF and DOCX documents (ЖЦП - Жизненный Цикл Проекта / Project Lifecycle Documents).

## Overview

The ЖЦП Parser is designed to automatically extract structured project information from unstructured Project Lifecycle Documents, including:
- Project phases
- Tasks within each phase
- Timeline information (start/end dates)
- Responsible persons and their roles
- Task dependencies and relationships

## Features

- **Multi-format Support**: PDF and DOCX document parsing
- **AI-Powered Extraction**: Uses LLMs to extract structured data
- **Fallback Mechanisms**: Supports multiple LLM providers (OpenAI, Anthropic, Ollama)
- **Data Validation**: Comprehensive validation and quality assurance
- **Error Handling**: Robust error handling and recovery
- **Configurable**: Flexible configuration for different use cases

## Prerequisites

- Go 1.21 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd parsing-ai
cd zhcp-parser-go
```

2. Install dependencies:
```bash
go mod download
```

## Configuration

### LLM Providers

The system supports multiple LLM providers:

1. **OpenAI** (GPT-4, GPT-4 Turbo)
2. **Anthropic** (Claude models)
3. **Ollama** (Local models like Llama 3)
4. **DeepSeek** (DeepSeek models)

Configure your preferred providers in `configs/llm_config.yaml`:

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
  - "ollama"    # Primary provider
  - "openai"    # Fallback 1
  - "anthropic" # Fallback 2
  - "deepseek"  # Fallback 3
```

### Environment Variables

For security, store API keys as environment variables:

```bash
export OPENAI_API_KEY="your-openai-api-key"
export ANTHROPIC_API_KEY="your-anthropic-api-key"
export DEEPSEEK_API_KEY="your-deepseek-api-key"
```

## Quick Start

### Basic Usage

Build and run the parser:

```bash
cd zhcp-parser-go
go build -o zhcp-parser cmd/zhcp-parser/main.go
./zhcp-parser parse path/to/your/document.pdf
```

### Command Line Interface

The parser provides a command-line interface with the following commands:

#### Parse a single document:
```bash
./zhcp-parser parse path/to/document.pdf
```

#### Parse with validation and enrichment:
```bash
./zhcp-parser parse path/to/document.pdf --validate --enrich
```

#### Parse with custom output:
```bash
./zhcp-parser parse path/to/document.pdf --output result.json
```

#### Parse with custom configuration:
```bash
./zhcp-parser parse path/to/document.pdf -c path/to/config.yaml
```

### Run the Demonstration

```bash
cd zhcp-parser-go
go run demo.go
```

This will run a complete demonstration showing:
- PDF document parsing
- DOCX document parsing
- Error handling capabilities

## API Reference

### Command Line Options

#### Global Options:
- `-c, --config`: Configuration file path (default: "configs/llm_config.yaml")

#### Parse Command Options:
- `-v, --validate`: Whether to perform validation (default: true)
- `-e, --enrich`: Whether to enrich data with computed fields (default: true)
- `-o, --output`: Output file for results (JSON format)

## Project Structure

```
parsing-ai/
└── zhcp-parser-go/
    ├── cmd/
    │   └── zhcp-parser/
    │       └── main.go             # Main entry point with CLI
    ├── configs/
    │   └── llm_config.yaml         # Configuration
    ├── internal/
    │   ├── ai/
    │   │   ├── llm_manager.go      # LLM manager
    │   │   ├── provider_interfaces.go # Provider interfaces
    │   │   ├── shared_types.go     # Shared types
    │   │   ├── types.go            # Type definitions
    │   │   └── llm_providers/      # LLM provider implementations
    │   │       ├── openai/
    │   │       ├── anthropic/
    │   │       ├── ollama/
    │   │       └── deepseek/
    │   ├── config/
    │   │   └── config.go           # Configuration management
    │   ├── errors/
    │   │   ├── error_handler.go    # Error handling
    │   │   └── types.go            # Error types
    │   ├── parser/
    │   │   ├── parser.go           # Main parser
    │   │   ├── config.go           # Parser config
    │   │   └── types.go            # Parser types
    │   ├── parsers/
    │   │   ├── text_preprocessor.go # Text preprocessing
    │   │   ├── pdf/                # PDF parsing
    │   │   │   ├── pdf_parser.go
    │   │   │   └── pdf_validator.go
    │   │   └── docx/               # DOCX parsing
    │   │       ├── docx_extractor.go
    │   │       └── docx_validator.go
    │   ├── transformers/
    │   │   ├── json_transformer.go # JSON transformation
    │   │   ├── data_enricher.go    # Data enrichment
    │   │   └── types.go            # Transformer types
    │   └── validators/
    │       ├── validation_pipeline.go # Validation pipeline
    │       ├── document_validator.go # Document validation
    │       └── types.go            # Validation types
    ├── testdata/
    │   └── sample_project.pdf      # Sample test data
    ├── demo.go                     # Demonstration script
    ├── go.mod                      # Go module definition
    └── README.md                   # This file
```

## Dependencies

- `github.com/pdfcpu/pdfcpu` - For PDF document parsing
- `github.com/baliance/gooxml` - For DOCX document parsing
- `github.com/sashabaranov/go-openai` - For OpenAI API integration
- `github.com/anthropics/anthropic-sdk-go` - For Anthropic API integration
- `github.com/ollama/ollama` - For Ollama integration
- `github.com/go-playground/validator/v10` - For data validation
- `github.com/rs/zerolog` - For logging
- `github.com/spf13/cobra` - For CLI commands
- `gopkg.in/yaml.v3` - For YAML configuration

## Performance Considerations

- **Large Documents**: The system handles large documents automatically with chunked processing
- **Rate Limiting**: Built-in rate limiting to prevent API abuse
- **Memory Management**: Efficient memory usage for large document processing

## Error Handling

The system implements comprehensive error handling with:
- Detailed error categorization
- Graceful degradation
- Fallback mechanisms
- Detailed error reporting

## Testing

To run tests:
```bash
cd zhcp-parser-go
go test ./... -v
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the GitHub repository.

---

**Note**: This is a production-ready Go implementation. For production use, ensure proper security measures and optimizations are in place.