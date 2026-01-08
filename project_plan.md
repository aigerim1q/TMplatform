# AI Module for Parsing Project Lifecycle Documents (ЖЦП)

## Project Overview
This project aims to create an AI module that accepts Project Lifecycle Documents (ЖЦП) and automatically extracts project structure including:
- Project phases
- Tasks
- Deadlines
- Responsible persons

## Technical Architecture

### Core Components
1. **Document Parser**: Handles PDF and DOCX document parsing
2. **AI Processing Engine**: Uses LLM to extract structured data
3. **Data Transformer**: Converts LLM responses to structured JSON
4. **Validation Layer**: Ensures extracted data quality

### Directory Structure
```
parsing-ai/
├── src/
│   ├── parsers/
│   │   ├── __init__.py
│   │   ├── pdf_parser.py
│   │   └── docx_parser.py
│   ├── ai/
│   │   ├── __init__.py
│   │   ├── llm_connector.py
│   │   └── prompt_engineering.py
│   ├── transformers/
│   │   ├── __init__.py
│   │   └── json_transformer.py
│   ├── validators/
│   │   ├── __init__.py
│   │   └── data_validator.py
│   ├── main.py
│   └── config.py
├── tests/
│   ├── samples/
│   │   ├── test_document.pdf
│   │   └── test_document.docx
│   ├── test_parsers.py
│   └── test_ai_extraction.py
├── requirements.txt
├── README.md
└── docs/
    └── api_spec.md
```

### JSON Output Schema
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
      "source_document": "string",
      "extraction_date": "YYYY-MM-DD",
      "confidence_score": "number"
    }
  }
}
```

### Implementation Steps

#### Phase 1: Document Parsing
- Implement PDF text extraction using PyPDF2 or similar
- Implement DOCX text extraction using python-docx
- Handle various document layouts and structures

#### Phase 2: AI Integration
- Set up LLM API connection (OpenAI, Claude, or local model)
- Design prompts for project structure extraction
- Implement response parsing and error handling

#### Phase 3: Data Transformation
- Convert LLM responses to structured JSON format
- Implement validation and data cleaning
- Add confidence scoring for extracted data

#### Phase 4: Testing and Validation
- Create test suite with sample documents
- Implement accuracy validation
- Performance testing and optimization

## Dependencies
- PyPDF2 or pdfplumber for PDF parsing
- python-docx for DOCX parsing
- OpenAI or similar for LLM integration
- Pydantic for data validation
- pytest for testing

## API Design
The module will expose:
- `parse_document(file_path)`: Main entry point for document parsing
- `extract_project_structure(text_content)`: AI-powered extraction
- `validate_extraction(extracted_data)`: Data validation
- `export_to_json(extracted_data)`: JSON export functionality