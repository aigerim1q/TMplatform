# ЖЦП Parser Usage Guide

## How to Use with Your Own Files

### 1. Setting up DeepSeek API

1. Get your DeepSeek API key from [DeepSeek Developer Portal](https://platform.deepseek.com/)
2. Set the API key as an environment variable:
   ```bash
   export DEEPSEEK_API_KEY="your-deepseek-api-key"
   ```
   On Windows Command Prompt:
   ```cmd
   set DEEPSEEK_API_KEY=your-deepseek-api-key
   ```
   On Windows PowerShell:
   ```powershell
   $env:DEEPSEEK_API_KEY="your-deepseek-api-key"
   ```

### 2. Configure the System

The system is already configured to use DeepSeek by default. Check your `configs/llm_config.yaml`:

```yaml
providers:
  deepseek:
    enabled: true
    api_key: "${DEEPSEEK_API_KEY}"
    model: "deepseek-chat"
    temperature: 0.1
    max_tokens: 2048

provider_priority:
  - "deepseek"
  - "ollama"
  - "openai"
  - "anthropic"
```

### 3. Running the Parser with Your Files

#### Build and Run:
```bash
cd zhcp-parser-go
go build -o zhcp-parser cmd/zhcp-parser/main.go
./zhcp-parser parse path/to/your/project_document.pdf
```

#### For DOCX files:
```bash
./zhcp-parser parse path/to/your/project_document.docx
```

### 4. Command Line Usage

#### Basic Usage:
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

#### Batch processing:
```bash
./zhcp-parser batch path/to/directory/with/documents/
```

### 5. Expected Output Format

The parser will return structured information like this:

```json
{
  "success": true,
  "project_structure": {
    "project": {
      "title": "Project Title",
      "description": "Project Description",
      "phases": [
        {
          "id": "phase_1",
          "name": "Phase Name",
          "description": "Phase Description",
          "start_date": "YYYY-MM-DD",
          "end_date": "YYYY-MM-DD",
          "tasks": [
            {
              "id": "task_1",
              "name": "Task Name",
              "description": "Task Description",
              "start_date": "YYYY-MM-DD",
              "end_date": "YYYY-MM-DD",
              "responsible_persons": [
                {
                  "name": "Person Name",
                  "role": "Role",
                  "contact": "Contact Info"
                }
              ],
              "dependencies": ["task_id"],
              "status": "planned"
            }
          ]
        }
      ]
    }
  },
  "extraction_metadata": {
    "confidence": 0.85,
    "status": "success",
    "processing_time": 12.34
  }
}
```

### 6. Troubleshooting

- **Document Format Issues**: Ensure your PDF/DOCX files are not scanned images (they should contain text)
- **API Connection**: Verify your DeepSeek API key is correctly set
- **Low Confidence**: Structure your documents with clear headings, phases, and task lists
- **Build Issues**: Make sure Go 1.21+ is installed and in your PATH

### 7. Environment Variables Setup

For convenience, you can create a `.env` file in your project root:

```bash
DEEPSEEK_API_KEY=your-actual-api-key-here
OPENAI_API_KEY=your-openai-api-key-here
ANTHROPIC_API_KEY=your-anthropic-api-key-here
```

The application will automatically read environment variables when needed.

### 8. Configuration File Location

The default configuration file is located at `configs/llm_config.yaml` but you can specify a custom location using the `-c` flag:

```bash
./zhcp-parser parse document.pdf -c /path/to/custom/config.yaml
```

The system is now ready to process your project lifecycle documents using DeepSeek R1!