# ЖЦП Parser API Documentation

## Overview

The ЖЦП (Project Lifecycle Document) Parser is an AI-powered module that automatically extracts project structure information from PDF and DOCX documents. It identifies project phases, tasks, deadlines, and responsible persons, converting unstructured document content into structured JSON data.

## Installation

### Prerequisites
- Python 3.8 or higher
- pip package manager

### Installation Steps
```bash
# Clone the repository
git clone <repository-url>
cd parsing-ai

# Create virtual environment (recommended)
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

### Required Dependencies
- `python-docx>=0.8.11` - For DOCX document parsing
- `PyMuPDF>=1.23.0` - For PDF document parsing
- `pdfplumber>=0.7.0` - Alternative PDF parsing
- `openai>=1.0.0` - For OpenAI API integration
- `anthropic>=0.5.0` - For Anthropic API integration
- `pydantic>=1.10.0` - For data validation
- `aiohttp>=3.8.0` - For async HTTP requests
- `pyyaml>=6.0` - For configuration management

## Configuration

### Configuration File Structure

Create a configuration file at `config/llm_config.yaml`:

```yaml
providers:
  openai:
    enabled: false  # Set to true to enable
    api_key: "your-openai-api-key"  # Use environment variable for security
    model: "gpt-4-turbo"
    temperature: 0.1
    max_tokens: 4096

  anthropic:
    enabled: false  # Set to true to enable
    api_key: "your-anthropic-api-key"  # Use environment variable for security
    model: "claude-3-sonnet-20240229"
    temperature: 0.1
    max_tokens: 4096

  ollama:
    enabled: true  # Set to true to enable local model
    model: "llama3"
    base_url: "http://localhost:11434"
    temperature: 0.1
    max_tokens: 4096

provider_priority:
  - "openai"
  - "anthropic"
  - "ollama"

retry_settings:
  max_retries: 3
  backoff_factor: 1.0
  status_codes: [429, 502, 503, 504]

rate_limiting:
  requests_per_minute: 60
  tokens_per_minute: 100000
```

### Environment Variables

For security, store API keys as environment variables:

```bash
export OPENAI_API_KEY="your-openai-api-key"
export ANTHROPIC_API_KEY="your-anthropic-api-key"
```

## Quick Start

### Basic Usage

```python
from src.main import ЖЦПParser

# Initialize the parser
parser = ЖЦПParser()

# Parse a PDF document
result = await parser.parse_document("path/to/your/document.pdf")

# Check if parsing was successful
if result['success']:
    project_structure = result['project_structure']
    print(f"Project title: {project_structure['project']['title']}")
    print(f"Number of phases: {len(project_structure['project']['phases'])}")
else:
    print(f"Parsing failed: {result['error']['message']}")
```

### Advanced Usage with Configuration

```python
from src.main import ЖЦПParser
from src.config import LLMConfig

# Load custom configuration
config = LLMConfig("path/to/custom/config.yaml")
parser = ЖЦПParser(config=config)

# Parse with additional options
result = await parser.parse_document(
    "path/to/document.pdf",
    validate=True,  # Enable validation
    enrich=True     # Enable data enrichment
)
```

## API Reference

### ЖЦПParser Class

#### Constructor
```python
parser = ЖЦПParser(config=None, log_level="INFO")
```

**Parameters:**
- `config` (LLMConfig, optional): Configuration object. If None, uses default configuration
- `log_level` (str, optional): Logging level. Default is "INFO"

#### Methods

##### `parse_document(document_path, validate=True, enrich=True)`

Parses a document and extracts project structure.

**Parameters:**
- `document_path` (str): Path to the PDF or DOCX document
- `validate` (bool, optional): Whether to perform validation. Default is True
- `enrich` (bool, optional): Whether to enrich data with computed fields. Default is True

**Returns:**
```python
{
    "success": bool,
    "project_structure": {
        # Project structure in standardized JSON format
    },
    "extraction_metadata": {
        "confidence": float,  # 0.0 to 1.0
        "status": str,        # "success", "partial", "failed", "validation_error"
        "processing_time": float,  # Processing time in seconds
        "validation_results": dict  # Validation results if validation was performed
    },
    "validation_errors": list,  # List of validation errors if any
    "processing_notes": list    # Additional processing notes
}
```

**Example:**
```python
result = await parser.parse_document("project_plan.pdf")
if result['success']:
    project = result['project_structure']['project']
    print(f"Title: {project['title']}")
    print(f"Phases: {len(project['phases'])}")
    for phase in project['phases']:
        print(f"  Phase: {phase['name']}")
        print(f"  Tasks: {len(phase['tasks'])}")
```

##### `get_error_summary()`

Returns a summary of recent errors.

**Returns:**
```python
{
    "total_errors": int,
    "by_category": {
        "parsing_error": int,
        "validation_error": int,
        # ... other categories
    },
    "by_severity": {
        "info": int,
        "warning": int,
        "error": int,
        "critical": int
    },
    "recent_errors": [
        {
            "error_id": str,
            "category": str,
            "severity": str,
            "message": str,
            "timestamp": str
        }
    ]
}
```

### Data Models

#### Project Structure Schema

The extracted data follows this JSON schema:

```json
{
  "project": {
    "title": "string",
    "description": "string",
    "phases": [
      {
        "id": "string",
        "name": "string",
        "description": "string",
        "start_date": "YYYY-MM-DD",
        "end_date": "YYYY-MM-DD",
        "tasks": [
          {
            "id": "string",
            "name": "string",
            "description": "string",
            "start_date": "YYYY-MM-DD",
            "end_date": "YYYY-MM-DD",
            "responsible_persons": [
              {
                "name": "string",
                "role": "string",
                "contact": "string"
              }
            ],
            "dependencies": ["string"],
            "status": "planned|in_progress|completed"
          }
        ]
      }
    ],
    "metadata": {
      "extraction_date": "string",
      "confidence_score": "number",
      "processing_notes": ["string"],
      "calculated_fields": {
        "project_start": "string",
        "project_end": "string",
        "total_duration_days": "number"
      },
      "data_quality_score": "number",
      "complexity_metrics": {
        "total_phases": "number",
        "total_tasks": "number",
        "total_responsibles": "number"
      }
    }
  }
}
```

## Usage Examples

### Example 1: Basic PDF Parsing

```python
import asyncio
from src.main import ЖЦПParser

async def parse_pdf_example():
    parser = ЖЦПParser()
    
    result = await parser.parse_document("sample_project.pdf")
    
    if result['success']:
        project = result['project_structure']['project']
        print(f"Project: {project['title']}")
        print(f"Confidence: {result['extraction_metadata']['confidence']:.2f}")
        
        for phase in project['phases']:
            print(f"\nPhase: {phase['name']}")
            for task in phase['tasks']:
                print(f"  - {task['name']}")
                if task['start_date'] and task['end_date']:
                    print(f"    Dates: {task['start_date']} to {task['end_date']}")
    else:
        print(f"Error: {result['error']['message']}")

# Run the example
asyncio.run(parse_pdf_example())
```

### Example 2: Batch Processing Multiple Documents

```python
import asyncio
from pathlib import Path
from src.main import ЖЦПParser

async def batch_process_documents():
    parser = ЖЦПParser()
    document_folder = Path("documents/")
    
    results = []
    
    for doc_path in document_folder.glob("*.pdf"):
        print(f"Processing: {doc_path.name}")
        result = await parser.parse_document(str(doc_path))
        results.append({
            'document': doc_path.name,
            'success': result['success'],
            'confidence': result.get('extraction_metadata', {}).get('confidence', 0),
            'error': result.get('error', {}).get('message', None)
        })
    
    # Print summary
    successful = sum(1 for r in results if r['success'])
    total = len(results)
    avg_confidence = sum(r['confidence'] for r in results if r['success']) / successful if successful > 0 else 0
    
    print(f"\nBatch processing summary:")
    print(f"Total documents: {total}")
    print(f"Successful: {successful}")
    print(f"Success rate: {successful/total:.2%}")
    print(f"Average confidence: {avg_confidence:.2f}")

# Run the example
asyncio.run(batch_process_documents())
```

### Example 3: Error Handling

```python
import asyncio
from src.main import ЖЦПParser

async def error_handling_example():
    parser = ЖЦПParser()
    
    try:
        result = await parser.parse_document("nonexistent.pdf")
        
        if not result['success']:
            error = result['error']
            print(f"Error occurred: {error['message']}")
            print(f"Error category: {error['category']}")
            print(f"Error ID: {error['error_id']}")
            
            # Get error summary
            error_summary = parser.get_error_summary()
            print(f"Total errors in session: {error_summary['total_errors']}")
    
    except Exception as e:
        print(f"Unexpected error: {str(e)}")

# Run the example
asyncio.run(error_handling_example())
```

## Configuration Options

### Provider Configuration

You can configure multiple LLM providers with different settings:

#### OpenAI Configuration
```yaml
providers:
  openai:
    enabled: true
    api_key: ${OPENAI_API_KEY}  # Use environment variable
    model: "gpt-4-turbo"
    temperature: 0.1
    max_tokens: 4096
    timeout: 30
```

#### Anthropic Configuration
```yaml
providers:
  anthropic:
    enabled: true
    api_key: ${ANTHROPIC_API_KEY}  # Use environment variable
    model: "claude-3-sonnet-20240229"
    temperature: 0.1
    max_tokens: 4096
    timeout: 45
```

#### Local Model (Ollama) Configuration
```yaml
providers:
  ollama:
    enabled: true
    model: "llama3"
    base_url: "http://localhost:11434"
    temperature: 0.1
    max_tokens: 4096
    timeout: 60
```

### Priority and Fallback

The system will attempt to use providers in the order specified by `provider_priority`, falling back to the next provider if one fails:

```yaml
provider_priority:
  - "openai"    # Primary provider
  - "anthropic" # Fallback 1
  - "ollama"    # Fallback 2 (local)
```

## Performance Considerations

### Processing Large Documents

For large documents, consider:

1. **Chunked Processing**: The system handles large documents automatically
2. **Asynchronous Processing**: Use async/await for better performance
3. **Caching**: Results are not cached by default, but you can implement caching
4. **Resource Management**: Monitor memory usage for very large documents

### Rate Limiting

The system implements rate limiting to prevent API abuse:

```yaml
rate_limiting:
  requests_per_minute: 60      # Max API requests per minute
  tokens_per_minute: 100000    # Max tokens per minute (for OpenAI)
```

## Troubleshooting

### Common Issues

#### 1. API Key Issues
- **Problem**: "Invalid API key" error
- **Solution**: Verify API keys are correctly set in environment variables or config file

#### 2. Document Format Issues
- **Problem**: "Unsupported document format" error
- **Solution**: Ensure document is valid PDF or DOCX format

#### 3. Network Issues
- **Problem**: Connection timeout errors
- **Solution**: Check network connectivity and increase timeout values in config

#### 4. Low Confidence Results
- **Problem**: Confidence scores below 0.5
- **Solution**: Ensure documents contain clear project structure information

### Debugging

Enable debug logging for detailed information:

```python
import logging
logging.basicConfig(level=logging.DEBUG)

parser = ЖЦПParser(log_level="DEBUG")
```

## Error Codes

| Code | Description | Solution |
|------|-------------|----------|
| PARSE_001 | Invalid document format | Verify file is PDF or DOCX |
| PARSE_002 | Document corrupted | Try opening document in viewer |
| LLM_001 | API key invalid | Check API key configuration |
| LLM_002 | Rate limit exceeded | Wait or increase rate limits |
| VALID_001 | Validation failed | Check document structure |
| CONF_001 | Configuration error | Verify config file format |

## Best Practices

### Document Preparation
- Use clear, structured documents with defined sections
- Include dates in standard formats (DD.MM.YYYY or YYYY-MM-DD)
- Use consistent terminology for phases and tasks
- Include responsible person names and roles

### Configuration Best Practices
- Use environment variables for API keys
- Start with local models (Ollama) for development
- Configure appropriate timeout values for your network
- Set up proper logging for production use

### Performance Optimization
- Process documents in batches for better throughput
- Use appropriate confidence thresholds for your use case
- Monitor and optimize API usage costs
- Implement result caching for frequently processed documents

## Integration Examples

### Web API Integration
```python
from fastapi import FastAPI, File, UploadFile
from src.main import ЖЦПParser
import tempfile
import os

app = FastAPI()
parser = ЖЦПParser()

@app.post("/parse-document/")
async def parse_document_api(file: UploadFile = File(...)):
    # Save uploaded file temporarily
    with tempfile.NamedTemporaryFile(delete=False, suffix=f".{file.filename.split('.')[-1]}") as temp_file:
        content = await file.read()
        temp_file.write(content)
        temp_path = temp_file.name
    
    try:
        result = await parser.parse_document(temp_path)
        return result
    finally:
        os.unlink(temp_path)
```

### Command Line Interface
```python
import argparse
import asyncio
import json
from src.main import ЖЦПParser

def main():
    parser = argparse.ArgumentParser(description='Parse ЖЦП documents')
    parser.add_argument('document_path', help='Path to PDF or DOCX document')
    parser.add_argument('--output', '-o', help='Output file path')
    parser.add_argument('--format', choices=['json', 'pretty'], default='json')
    
    args = parser.parse_args()
    
    async def run():
        parser_instance = ЖЦПParser()
        result = await parser_instance.parse_document(args.document_path)
        
        if args.format == 'pretty':
            output = json.dumps(result, indent=2, ensure_ascii=False)
        else:
            output = json.dumps(result, ensure_ascii=False)
        
        if args.output:
            with open(args.output, 'w', encoding='utf-8') as f:
                f.write(output)
        else:
            print(output)
    
    asyncio.run(run())

if __name__ == "__main__":
    main()
```

This comprehensive API documentation provides all the information needed to integrate and use the ЖЦП parsing module effectively.