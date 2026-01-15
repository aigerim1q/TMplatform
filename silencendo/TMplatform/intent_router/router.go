package router

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"silencendo/llm"
)

const UniversalSystemPrompt = `You are the reasoning engine of the system. Process the user's input and respond with ONLY a valid JSON object. No explanations, no markdown, no additional text.

Return a JSON object with these exact fields:
{
  "intent": "answer_question | list_projects | project_current | task_create | task_update | stage_list | context_show | help | ask_clarification | chat | create_project | add_stages_tasks | assign_responsible | show_project",
  "confidence": 0.0-1.0,
  "parameters": {},
  "response": "natural language response for the user"
}

Available intents:
- answer_question: for general questions
- list_projects: when user wants to see projects
- project_current: when user wants to see current project details (equivalent to show_project, shows the plan and all details)
- task_create: when user wants to create tasks
- task_update: when user wants to update tasks
- stage_list: when user wants to list stages
- context_show: when user wants to see context
- help: when user asks for help
- ask_clarification: when you're uncertain about intent
- chat: for general conversation
- create_project: when user wants to create a project
- add_stages_tasks: when user wants to add stages or tasks
- assign_responsible: when user wants to assign tasks/stages
- show_project: when user wants to see project details (stages, tasks, assignments)

CRITICAL: For "project_current", "show the plan", "show project", or any request to view current project details, use intent "project_current" or "show_project" which will display the full project plan with stages, tasks, and assignments.

Rules:
- If unsure about intent, use "ask_clarification"
- Never refuse due to formatting
- Always return valid JSON
- You are the reasoning engine of the system`

// IntentMatchLocal mirrors the chatbot.IntentMatch struct
type IntentMatchLocal struct {
	Intent       string
	ProjectTitle string
	Description  string
	StagePlans   []interface{} // Will be defined as needed
	EntityType   string
	EntityName   string
	AssigneeName string
	Filters      interface{} // Will be defined as needed
	Members      []string
	Confidence   float64
}

// Router handles LLM-based intent routing
type Router struct {
	LLMClient llm.LLMClient
}

// New creates a new Router instance
func New(llmClient llm.LLMClient) *Router {
	return &Router{
		LLMClient: llmClient,
	}
}

// ProcessInput sends the user input to the LLM and returns the parsed response
func (r *Router) ProcessInput(ctx context.Context, input string) (*LLMResponse, error) {
	messages := []llm.Message{
		{Role: "system", Content: UniversalSystemPrompt},
		{Role: "user", Content: input},
	}

	response, err := r.LLMClient.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response from LLM: %w", err)
	}

	// Try to parse the response as JSON
	llmResponse, err := r.parseResponse(response)
	if err != nil {
		// Try to fix the JSON response once
		fixedResponse, fixErr := r.fixJSONResponse(ctx, input, response)
		if fixErr != nil {
			// If fixing fails, return a default clarification response
			return &LLMResponse{
				Intent:     "ask_clarification",
				Confidence: 0.5,
				Parameters: make(map[string]interface{}),
				Response:   "I'm having trouble understanding your request. Could you please rephrase or be more specific?",
			}, nil
		}

		llmResponse, err = r.parseResponse(fixedResponse)
		if err != nil {
			// Second attempt also failed, return clarification
			return &LLMResponse{
				Intent:     "ask_clarification",
				Confidence: 0.5,
				Parameters: make(map[string]interface{}),
				Response:   "I'm having trouble understanding your request. Could you please rephrase or be more specific?",
			}, nil
		}
	}

	// Validate the response
	if !llmResponse.Validate() {
		return &LLMResponse{
			Intent:     "ask_clarification",
			Confidence: 0.5,
			Parameters: make(map[string]interface{}),
			Response:   "I'm having trouble understanding your request. Could you please rephrase or be more specific?",
		}, nil
	}

	return llmResponse, nil
}

// parseResponse extracts and parses the JSON from the LLM response
func (r *Router) parseResponse(response string) (*LLMResponse, error) {
	// Find JSON in the response (may be surrounded by other text)
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd < jsonStart {
		return nil, fmt.Errorf("no JSON object found in response")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var llmResponse LLMResponse
	if err := json.Unmarshal([]byte(jsonStr), &llmResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &llmResponse, nil
}

// fixJSONResponse attempts to fix invalid JSON by asking the LLM to correct it
func (r *Router) fixJSONResponse(ctx context.Context, originalInput, malformedResponse string) (string, error) {
	// First, try to extract any JSON-like structure from the malformed response
	if extracted := extractJSONFromText(malformedResponse); extracted != "" {
		// Test if the extracted JSON is valid
		var testObj map[string]interface{}
		if err := json.Unmarshal([]byte(extracted), &testObj); err == nil {
			// Check if it has the required fields
			if _, hasIntent := testObj["intent"]; hasIntent {
				return extracted, nil
			}
		}
	}

	// If extraction didn't work, ask the LLM to fix it
	fixPrompt := fmt.Sprintf(
		"The following response is not valid JSON:\n\n%s\n\n"+
			"Please fix the JSON and return ONLY the corrected JSON object with no additional text. "+
			"The JSON should have these fields: intent (string), confidence (float), parameters (object), response (string).",
		malformedResponse,
	)

	messages := []llm.Message{
		{Role: "system", Content: UniversalSystemPrompt},
		{Role: "user", Content: fixPrompt},
	}

	fixedResponse, err := r.LLMClient.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to fix JSON response: %w", err)
	}

	// Verify that the fixed response is actually valid JSON
	jsonStart := strings.Index(fixedResponse, "{")
	jsonEnd := strings.LastIndex(fixedResponse, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd < jsonStart {
		return "", fmt.Errorf("fixed response is not valid JSON")
	}

	// Double-check the extracted JSON is valid
	extractedFixed := fixedResponse[jsonStart : jsonEnd+1]
	var testObj map[string]interface{}
	if err := json.Unmarshal([]byte(extractedFixed), &testObj); err != nil {
		return "", fmt.Errorf("fixed response is not valid JSON: %w", err)
	}

	return extractedFixed, nil
}

// extractJSONFromText tries to extract a JSON object from text that may contain other content
func extractJSONFromText(text string) string {
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd < jsonStart {
		return ""
	}

	return text[jsonStart : jsonEnd+1]
}

// ConvertToIntentMatch converts the LLM response to a local IntentMatch
func (r *Router) ConvertToIntentMatch(llmResponse *LLMResponse, originalInput string) IntentMatchLocal {
	intent := llmResponse.ToChatbotIntent()

	match := IntentMatchLocal{
		Intent:     intent,
		Confidence: llmResponse.Confidence,
	}

	// Extract parameters based on intent
	switch intent {
	case IntentCreateProject:
		if title, ok := llmResponse.Parameters["project_title"].(string); ok {
			match.ProjectTitle = title
		} else {
			// Try to extract project title from the original input if not in parameters
			// This is a simplified extraction - in practice, the LLM should provide it
			match.ProjectTitle = extractProjectTitle(originalInput)
		}

		if description, ok := llmResponse.Parameters["description"].(string); ok {
			match.Description = description
		}

	case IntentAddStagesAndTasks:
		// Extract stages and tasks from parameters if available
		if _, ok := llmResponse.Parameters["stages"].([]interface{}); ok {
			// Note: We may need to handle this differently based on how the LLM formats stages
		}

		if projectName, ok := llmResponse.Parameters["project_title"].(string); ok {
			match.ProjectTitle = projectName
		}

	case IntentAssignResponsible:
		if entityType, ok := llmResponse.Parameters["entity_type"].(string); ok {
			match.EntityType = entityType
		}
		if entityName, ok := llmResponse.Parameters["entity_name"].(string); ok {
			match.EntityName = entityName
		}
		if assigneeName, ok := llmResponse.Parameters["assignee_name"].(string); ok {
			match.AssigneeName = assigneeName
		}
	}

	return match
}

// extractProjectTitle tries to extract a project title from the input
func extractProjectTitle(input string) string {
	// Look for common patterns in the input
	lower := strings.ToLower(input)

	// Patterns to extract project title
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?:create|make|build|start|new)\s+project\s+(.+)`),
		regexp.MustCompile(`(?:create|make|build|start|new)\s+(.+)\s+project`),
		regexp.MustCompile(`project\s+(.+)`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(lower)
		if len(matches) > 1 {
			title := strings.TrimSpace(matches[1])
			// Clean up the title
			title = strings.Trim(title, "\"'")
			return title
		}
	}

	return ""
}
