# Natural User Input Examples

Here are 5 examples of natural user inputs and the expected JSON output from the LLM:

## Example 1: Project Creation
**Input:** "I want to start a new project for building a mobile app called TaskMaster"
**Expected JSON Output:**
```json
{
  "intent": "create_project",
  "confidence": 0.9,
  "parameters": {
    "project_title": "TaskMaster",
    "description": "building a mobile app"
  },
  "response": "Creating a new project called TaskMaster for building a mobile app."
}
```

## Example 2: Project Listing
**Input:** "What projects am I currently working on?"
**Expected JSON Output:**
```json
{
  "intent": "list_projects",
  "confidence": 0.95,
  "parameters": {},
  "response": "Here are the projects you're currently working on:"
}
```

## Example 3: Task Assignment
**Input:** "Can you assign John to handle the UI design for the TaskMaster project?"
**Expected JSON Output:**
```json
{
  "intent": "assign_responsible",
  "confidence": 0.85,
  "parameters": {
    "entity_type": "task",
    "entity_name": "UI design",
    "assignee_name": "John",
    "project_title": "TaskMaster"
  },
  "response": "Assigning John to handle the UI design task for the TaskMaster project."
}
```

## Example 4: Adding Stages
**Input:** "We need to add a testing phase to the TaskMaster project with unit tests and integration tests"
**Expected JSON Output:**
```json
{
  "intent": "add_stages_tasks",
  "confidence": 0.8,
  "parameters": {
    "project_title": "TaskMaster",
    "stages": [
      {
        "title": "Testing",
        "tasks": ["unit tests", "integration tests"]
      }
    ]
  },
  "response": "Adding a testing phase to the TaskMaster project with unit tests and integration tests."
}
```

## Example 5: Project Details Request
**Input:** "Show me the current status of the TaskMaster project"
**Expected JSON Output:**
```json
{
  "intent": "show_project",
  "confidence": 0.9,
  "parameters": {
    "project_title": "TaskMaster"
  },
  "response": "Here's the current status of the TaskMaster project:"
}