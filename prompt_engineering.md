# Prompt Engineering for ЖЦП Document Parsing

## Overview
This document outlines the prompt engineering strategy for extracting project structure information from ЖЦП (Project Lifecycle Documents) using AI/LLM technology.

## Target Output Format

The AI should return structured JSON with the following schema:

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
      "confidence_score": "number",
      "processing_notes": ["string"]
    }
  }
}
```

## Core Prompt Template

```
You are a project management expert specializing in Russian project lifecycle documents (ЖЦП). 
Extract the complete project structure from the following document content, identifying:

1. Project phases (main stages of the project)
2. Tasks within each phase
3. Timeline information (start/end dates)
4. Responsible persons and their roles
5. Task dependencies and relationships

Document content:
{document_content}

Extract this information and return ONLY a valid JSON object with the following structure:
{json_schema}

Important guidelines:
- Use Russian terminology where appropriate in the output
- If dates are not explicitly mentioned, set to null
- If responsible persons are not explicitly mentioned, set to empty array
- Estimate confidence scores based on clarity of information in the document
- Use UUID-like strings for IDs (e.g., "phase_1", "task_1_1")
- Keep descriptions concise but informative
- If you cannot determine certain information, use null values
- Do not include any explanatory text outside the JSON
```

## Enhanced Prompt with Examples

```
You are a project management expert specializing in Russian project lifecycle documents (ЖЦП). Extract structured project information from the provided document.

The output must be a valid JSON object with the following structure:
{json_schema}

Here are examples of proper extraction:

EXAMPLE 1 - Input snippet:
"Этап 1: Подготовительный этап (январь 2024 - март 2024)
1.1. Разработка технического задания - ответственный: Иванов И.И., срок выполнения: до 15.01.2024
1.2. Согласование документации - ответственный: Петров П.П., срок выполнения: 16.01.2024 - 31.01.2024"

EXAMPLE 1 - Expected output:
{
  "phases": [
    {
      "id": "phase_1",
      "name": "Подготовительный этап",
      "description": "",
      "start_date": "2024-01-01",
      "end_date": "2024-03-31",
      "tasks": [
        {
          "id": "task_1_1",
          "name": "Разработка технического задания",
          "description": "",
          "start_date": null,
          "end_date": "2024-01-15",
          "responsible_persons": [
            {
              "name": "Иванов И.И.",
              "role": "",
              "contact": ""
            }
          ],
          "dependencies": [],
          "status": "planned"
        },
        {
          "id": "task_1_2",
          "name": "Согласование документации",
          "description": "",
          "start_date": "2024-01-16",
          "end_date": "2024-01-31",
          "responsible_persons": [
            {
              "name": "Петров П.П.",
              "role": "",
              "contact": ""
            }
          ],
          "dependencies": [],
          "status": "planned"
        }
      ]
    }
  ]
}

Now process this document content:
{document_content}

Return ONLY the JSON structure with extracted information. No additional text.
```

## Specialized Prompts for Different Document Types

### For Table-Based Documents
```
This document contains project information in table format. Extract phases, tasks, deadlines, and responsible persons from the table structure.

Focus on:
- Table headers that indicate phase/task columns
- Date columns for timeline information
- Person/responsible columns
- Row relationships that indicate task dependencies

Document content:
{table_content}
```

### For Narrative Documents
```
This document describes the project in narrative form. Extract structured information by identifying:
- Section headers that indicate project phases
- Sentences that describe tasks and activities
- Date references and timeframes
- References to responsible persons and roles

Pay special attention to phrases like:
- "на этапе..." (at the stage of...)
- "ответственный..." (responsible...)
- "до [дата]" (by [date])
- "в течение [период]" (within [period])

Document content:
{narrative_content}
```

## Prompt Variants for Accuracy Improvement

### Variant 1: Step-by-Step Extraction
```
First, identify all project phases mentioned in the document.
Second, for each phase, identify the tasks.
Third, extract timeline information for phases and tasks.
Fourth, identify responsible persons and their roles.
Finally, organize all information into the required JSON structure.

Document content:
{document_content}
```

### Variant 2: Role-Based Extraction
```
As a Russian project management analyst, you are reviewing a project lifecycle document. 
Your task is to create a structured project plan from this document by identifying:
- What needs to be done (tasks)
- When it needs to be done (timelines)
- Who is responsible (personnel)
- In what sequence (phases and dependencies)

Document content:
{document_content}
```

## Error Handling in Prompts

### Handling Ambiguous Information
```
If dates are ambiguous (e.g., "next month"), set to null and add a note in processing_notes.
If person names are unclear, set to null but note the ambiguity.
If task dependencies are implied but not explicit, set dependencies to empty array.
```

### Handling Different Date Formats
```
Convert all date formats to YYYY-MM-DD:
- DD.MM.YYYY → YYYY-MM-DD
- DD/MM/YYYY → YYYY-MM-DD
- Month names in Russian (январь, февраль, etc.) → appropriate dates
- Relative dates (next week, end of month) → null with explanation in notes
```

## Confidence Scoring Guidelines

Include in prompt:
```
Provide confidence scores (0.0 to 1.0) based on:
- 0.9-1.0: Explicit, clear information directly stated
- 0.7-0.8: Information clearly implied or easily inferred
- 0.5-0.6: Information somewhat ambiguous but reasonable inference
- 0.3-0.4: Information unclear but possible interpretation
- 0.1-0.2: Information very unclear, mostly guessing
```

## Testing and Validation Prompts

### Validation Prompt
```
Review the extracted JSON structure and validate:
1. All required fields are present
2. Date formats are correct (YYYY-MM-DD or null)
3. Relationships between phases and tasks are logical
4. No duplicate IDs
5. All referenced dependencies exist

Original document context: {original_context}
Extracted JSON: {extracted_json}
```

## Performance Optimization

### Chunking Strategy for Large Documents
For documents exceeding LLM context limits:
1. Process document in overlapping chunks
2. Extract partial information from each chunk
3. Merge and resolve conflicts between chunks
4. Validate consistency across the entire document

### Iterative Refinement
```
First pass: Extract high-level phases and major tasks
Second pass: Extract detailed task information, dates, and responsibilities
Third pass: Validate and connect dependencies between tasks
```

## Quality Assurance Prompts

### Fact-Checking Prompt
```
Compare the extracted information with the original document:
- Verify that all extracted phases actually appear in the document
- Verify that all extracted tasks are supported by document content
- Verify that all dates and person names match document text
- Identify any information that might be hallucinated
```

This prompt engineering approach ensures consistent, accurate extraction of project structure information from ЖЦП documents while handling the specific requirements of Russian language project management documentation.