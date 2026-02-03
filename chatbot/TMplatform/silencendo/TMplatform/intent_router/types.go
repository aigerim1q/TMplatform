package router

// Define constants directly instead of importing from chatbot
const (
	IntentCreateProject     = "create_project"
	IntentAddStagesAndTasks = "add_stages_tasks"
	IntentAssignResponsible = "assign_responsible"
	IntentListProjects      = "list_projects"
	IntentShowProject       = "show_project"
	IntentUnknown           = "unknown"
)

// LLMResponse represents the standard JSON response from the LLM
type LLMResponse struct {
	Intent     string                 `json:"intent"`
	Confidence float64                `json:"confidence"`
	Parameters map[string]interface{} `json:"parameters"`
	Response   string                 `json:"response"`
}

// IntentMap maps LLM intent strings to chatbot intents
var IntentMap = map[string]string{
	"answer_question":    "answer_question", // This will be handled separately
	"list_projects":      "list_projects",
	"project_current":    "show_project", // Maps to show_project in chatbot - shows current project details/plan
	"task_create":        "add_stages_tasks",
	"task_update":        "task_update",  // Will need to add this
	"stage_list":         "show_project", // Shows project details
	"context_show":       "show_project", // Shows project details
	"help":               "help",
	"ask_clarification":  "ask_clarification",
	"chat":               "chat",
	"create_project":     "create_project",
	"add_stages_tasks":   "add_stages_tasks",
	"assign_responsible": "assign_responsible",
	"show_project":       "show_project",
	"unknown":            "unknown",
}

// ToChatbotIntent converts LLM intent to chatbot intent
func (lr *LLMResponse) ToChatbotIntent() string {
	if intent, exists := IntentMap[lr.Intent]; exists {
		return intent
	}
	return IntentUnknown
}

// Validate ensures the LLM response has required fields
func (lr *LLMResponse) Validate() bool {
	if lr.Intent == "" || lr.Response == "" {
		return false
	}

	// Set default confidence if missing
	if lr.Confidence == 0 {
		lr.Confidence = 0.5
	}

	// Initialize parameters if missing
	if lr.Parameters == nil {
		lr.Parameters = make(map[string]interface{})
	}

	return true
}
