# LLM Intent Router

This module implements a centralized LLM-based intent detection and routing system that replaces the previous rule-based approach.

## Architecture Overview

The new architecture follows a "LLM-as-the-central-reasoning-engine" pattern where:

1. **Every user input is processed by the LLM** - No more text-based routing, keyword matching, or regex checks
2. **Structured JSON responses** - The LLM returns standardized JSON objects with intent, confidence, parameters, and response
3. **Thin routing layer** - The system routes based solely on the `intent` field from LLM output

## Key Components

### Universal System Prompt
The LLM uses a consistent system prompt that forces it to:
- Always return valid JSON only
- Never return markdown or explanations
- Follow the specified schema with intent, confidence, parameters, and response fields
- Act as the reasoning engine of the system

### JSON Schema
The LLM must always return JSON in this format:
```json
{
  "intent": "answer_question | list_projects | project_current | task_create | task_update | stage_list | context_show | help | ask_clarification | chat | create_project | add_stages_tasks | assign_responsible | show_project",
  "confidence": 0.0-1.0,
  "parameters": {},
  "response": "natural language response for the user"
}
```

### Intent Mapping
Intents are mapped from LLM responses to the existing chatbot intents:
- `list_projects` → `IntentListProjects`
- `create_project` → `IntentCreateProject`
- `add_stages_tasks` → `IntentAddStagesAndTasks`
- `assign_responsible` → `IntentAssignResponsible`
- `show_project` → `IntentShowProject`
- And others...

## Fallback Mechanisms

The system handles various error conditions:
1. **Invalid JSON**: Attempts to fix JSON by asking LLM to correct itself
2. **LLM Failures**: Falls back to rule-based detection if LLM is unavailable
3. **Unclear Intents**: Uses `ask_clarification` intent when LLM is uncertain

## Benefits

- **Natural Language Processing**: Users can express requests in natural language without specific keywords
- **Reduced "Unknown Command" Errors**: LLM interprets intent rather than relying on exact matches
- **Flexible Input Handling**: Complex requests are understood without rigid command structures
- **Maintainable Architecture**: Centralized reasoning reduces branching logic

## Integration Points

- **chatbot/intent.go**: Now uses LLM-based detection with fallback to rule-based
- **chatbot/handler.go**: Maintains the same interface but can now handle LLM responses directly
- **main.go**: Initializes the global router and sets it for use throughout the application

## Usage Examples

Instead of requiring users to say "list projects", they can now say:
- "Show me what I'm working on"
- "What projects do I have?"
- "Give me a list of projects"
- "I want to see my projects"

All will be correctly interpreted by the LLM and routed appropriately.