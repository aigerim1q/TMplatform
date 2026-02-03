package router

import (
	"testing"
)

// TestValidate tests the validation of LLM responses
func TestValidate(t *testing.T) {
	// Test valid response
	validResponse := &LLMResponse{
		Intent:     "list_projects",
		Confidence: 0.8,
		Parameters: map[string]interface{}{},
		Response:   "Here are your projects",
	}
	if !validResponse.Validate() {
		t.Error("Expected valid response to return true")
	}

	// Test response without required fields
	invalidResponse := &LLMResponse{
		Intent:     "",
		Confidence: 0,
		Parameters: nil,
		Response:   "",
	}
	if invalidResponse.Validate() {
		t.Error("Expected invalid response to return false")
	}

	// Test response with missing confidence gets default
	partialResponse := &LLMResponse{
		Intent:     "list_projects",
		Confidence: 0,   // Missing confidence
		Parameters: nil, // Missing parameters
		Response:   "Here are your projects",
	}
	if !partialResponse.Validate() {
		t.Error("Expected partial response to return true after validation")
	}
	if partialResponse.Confidence != 0.5 {
		t.Errorf("Expected default confidence 0.5, got %f", partialResponse.Confidence)
	}
	if partialResponse.Parameters == nil {
		t.Error("Expected parameters to be initialized")
	}
}

// TestExtractJSONFromText tests the JSON extraction utility function
func TestExtractJSONFromTextFunction(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    `Some text before {"intent": "list_projects", "confidence": 0.9, "parameters": {}, "response": "Here are your projects"} Some text after`,
			expected: `{"intent": "list_projects", "confidence": 0.9, "parameters": {}, "response": "Here are your projects"}`,
		},
		{
			input:    `{"intent": "list_projects", "confidence": 0.9}`,
			expected: `{"intent": "list_projects", "confidence": 0.9}`,
		},
		{
			input:    `No JSON here`,
			expected: ``,
		},
	}

	for _, tc := range testCases {
		result := extractJSONFromText(tc.input)
		if result != tc.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", tc.input, tc.expected, result)
		}
	}
}

// TestIntentMapping tests the mapping of LLM intents to chatbot intents
func TestIntentMappingFunction(t *testing.T) {
	response := &LLMResponse{
		Intent: "list_projects",
	}

	mappedIntent := response.ToChatbotIntent()
	if mappedIntent != IntentListProjects {
		t.Errorf("Expected chatbot intent '%s', got '%s'", IntentListProjects, mappedIntent)
	}

	// Test unknown intent mapping
	response.Intent = "unknown_intent"
	mappedIntent = response.ToChatbotIntent()
	if mappedIntent != IntentUnknown {
		t.Errorf("Expected chatbot intent '%s', got '%s'", IntentUnknown, mappedIntent)
	}
}
