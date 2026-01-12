package cli

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"silencendo/chatbot"
	"silencendo/llm"
	"silencendo/sources"
)

// IntentRouter evaluates free-form input and decides routing between project and document handlers.
type IntentRouter struct {
	sourceManager *sources.SourceManager
	detectIntent  func(string) string
	llmClient     llm.LLMClient
	useLLM        bool
}

func NewIntentRouter(sourceManager *sources.SourceManager, detectIntent func(string) string) *IntentRouter {
	useLLM := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY")) != ""
	var client llm.LLMClient
	if useLLM {
		client = llm.CreateLLMClient()
		if _, isMock := client.(*llm.MockLLM); isMock {
			useLLM = false
		}
	}

	return &IntentRouter{sourceManager: sourceManager, detectIntent: detectIntent, llmClient: client, useLLM: useLLM}
}

// IntentType enumerates the routing categories.
type IntentType string

const (
	IntentProject   IntentType = "project"
	IntentDocument  IntentType = "document"
	IntentKnowledge IntentType = "knowledge"
	IntentAmbiguous IntentType = "ambiguous"
)

// IntentResult captures routing output.
type IntentResult struct {
	Type         IntentType
	Confidence   float64
	ProjectMatch *chatbot.IntentMatch
	Resolved     *ResolvedProject
	LLMFailed    bool
}

type ResolvedProject struct {
	ID         string
	Title      string
	Confidence float64
}

func (r *IntentRouter) Route(ctx context.Context, input string, resolved *ResolvedProject) IntentResult {
	normalized := strings.TrimSpace(input)
	lower := strings.ToLower(normalized)

	if r.useLLM && r.llmClient != nil {
		if classified := r.classifyWithDeepSeek(ctx, normalized, resolved); classified != nil {
			switch classified.Domain {
			case "project":
				match := mapClassifiedProjectIntent(normalized, resolved, classified.Intent)
				if match.Intent != chatbot.IntentUnknown {
					return IntentResult{Type: IntentProject, Confidence: classified.Confidence, ProjectMatch: &match, Resolved: resolved}
				}
			case "docs":
				if classified.Intent == "docs_command" {
					return IntentResult{Type: IntentDocument, Confidence: classified.Confidence}
				}
				if classified.Intent == "docs_qa" {
					return IntentResult{Type: IntentKnowledge, Confidence: classified.Confidence}
				}
			}
		}

		if resolved != nil {
			if match, conf, err := r.classifyProjectWithDeepSeek(ctx, normalized); err == nil && match != nil {
				match.ProjectTitle = resolved.Title
				return IntentResult{Type: IntentProject, Confidence: conf, ProjectMatch: match, Resolved: resolved}
			}
		}

		return IntentResult{Type: IntentKnowledge, Confidence: 0, Resolved: resolved, LLMFailed: true}
	}

	// No LLM available; fall back to lightweight heuristics.
	if isProjectListQuery(lower) && (resolved == nil || strings.TrimSpace(resolved.ID) == "") {
		match := chatbot.IntentMatch{Intent: chatbot.IntentListProjects, Confidence: 0.9}
		return IntentResult{Type: IntentProject, Confidence: 0.9, ProjectMatch: &match, Resolved: resolved}
	}

	projectMatch := chatbot.DetectIntent(normalized)
	if projectMatch.Intent != chatbot.IntentUnknown {
		conf := projectMatch.Confidence
		if conf == 0 {
			conf = 0.6
		}
		if resolved != nil && strings.TrimSpace(resolved.ID) != "" && strings.TrimSpace(projectMatch.ProjectTitle) == "" {
			projectMatch.ProjectTitle = resolved.Title
		}
		return IntentResult{Type: IntentProject, Confidence: conf, ProjectMatch: &projectMatch, Resolved: resolved}
	}

	intentFn := r.detectIntent
	if intentFn == nil {
		intentFn = func(string) string { return "question" }
	}

	docOrQuestion := intentFn(normalized)
	if docOrQuestion == "edit" {
		return IntentResult{Type: IntentDocument, Confidence: 0.6}
	}

	return IntentResult{Type: IntentKnowledge, Confidence: 0.55}
}

type llmIntentResponse struct {
	IntentType    string  `json:"intent_type"`
	ProjectIntent string  `json:"project_intent"`
	Confidence    float64 `json:"confidence"`
}

type projectClassification struct {
	Domain     string  `json:"domain"`
	Intent     string  `json:"intent"`
	Confidence float64 `json:"confidence"`
}

func (r *IntentRouter) classifyProjectWithDeepSeek(ctx context.Context, input string) (*chatbot.IntentMatch, float64, error) {
	if !r.useLLM || r.llmClient == nil {
		return nil, 0, errors.New("llm disabled")
	}

	classifyCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	messages := []llm.Message{
		{Role: "system", Content: "You route CLI requests. Classify each message into intent_type of project|document|knowledge. If the user asks to view, show, or describe a project plan or project content (e.g., 'show the plan', 'show the content', 'what is the content'), set intent_type to project and project_intent to show_project. project_intent options: create_project, add_stages_tasks, assign_responsible, list_projects, show_project, unknown. Respond with a single JSON object and nothing else."},
		{Role: "user", Content: input},
	}

	raw, err := r.llmClient.Generate(classifyCtx, messages)
	if err != nil {
		return nil, 0, err
	}

	parsed, err := decodeLLMIntent(raw)
	if err != nil {
		return nil, 0, err
	}

	if strings.ToLower(strings.TrimSpace(parsed.IntentType)) != string(IntentProject) {
		return nil, 0, errors.New("llm did not classify as project")
	}

	conf := clampConfidence(parsed.Confidence)
	if conf < 0.5 {
		return nil, 0, errors.New("llm confidence too low")
	}

	mapped := mapProjectIntent(parsed.ProjectIntent)
	match := chatbot.DetectIntent(input)
	if match.Intent == chatbot.IntentUnknown && mapped != chatbot.IntentUnknown {
		match.Intent = mapped
	}
	if match.Intent == chatbot.IntentUnknown {
		return nil, 0, errors.New("project intent unknown")
	}

	if match.Confidence < conf {
		match.Confidence = conf
	}

	return &match, match.Confidence, nil
}

func decodeLLMIntent(raw string) (*llmIntentResponse, error) {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end <= start {
		return nil, errors.New("no json payload found")
	}

	segment := raw[start : end+1]
	var resp llmIntentResponse
	if err := json.Unmarshal([]byte(segment), &resp); err != nil {
		return nil, err
	}

	resp.IntentType = strings.TrimSpace(strings.ToLower(resp.IntentType))
	resp.ProjectIntent = strings.TrimSpace(strings.ToLower(resp.ProjectIntent))
	resp.Confidence = clampConfidence(resp.Confidence)
	return &resp, nil
}

func mapProjectIntent(intent string) string {
	switch strings.ToLower(strings.TrimSpace(intent)) {
	case "create_project", "createproject", "new_project":
		return chatbot.IntentCreateProject
	case "add_stages_tasks", "add_plan", "create_plan", "plan_creation":
		return chatbot.IntentAddStagesAndTasks
	case "assign_responsible", "assign", "ownership":
		return chatbot.IntentAssignResponsible
	case "list_projects", "project_list", "show_projects":
		return chatbot.IntentListProjects
	case "show_project", "view_project", "show_plan", "view_plan", "project_details":
		return chatbot.IntentShowProject
	default:
		return chatbot.IntentUnknown
	}
}

func (r *IntentRouter) classifyWithDeepSeek(ctx context.Context, input string, resolved *ResolvedProject) *projectClassification {
	if !r.useLLM || r.llmClient == nil {
		return nil
	}

	classifyCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	messages := []llm.Message{
		{Role: "system", Content: "Classify the user message. Respond with JSON only: {\"domain\": \"project|docs|unknown\", \"intent\": one of [project_show_content, project_show_plan, project_create_plan, project_list_tasks, project_list_stages, project_list_members, project_count_done_tasks, docs_command, docs_qa, unknown], \"confidence\": 0-1}. If the user asks about project content (e.g., 'what is the content', 'show the content'), choose intent=project_show_content. Do not guess project names. If no active project is mentioned, still classify intent but NEVER invent a project. Return ONLY JSON."},
		{Role: "user", Content: input},
	}

	raw, err := r.llmClient.Generate(classifyCtx, messages)
	if err != nil {
		return nil
	}

	classified, err := decodeProjectClassification(raw)
	if err != nil {
		return nil
	}

	classified.Confidence = clampConfidence(classified.Confidence)
	if classified.Confidence < 0.6 {
		return nil
	}

	return classified
}

func decodeProjectClassification(raw string) (*projectClassification, error) {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end <= start {
		return nil, errors.New("no json payload found")
	}

	segment := raw[start : end+1]
	var resp projectClassification
	if err := json.Unmarshal([]byte(segment), &resp); err != nil {
		return nil, err
	}

	resp.Domain = strings.TrimSpace(strings.ToLower(resp.Domain))
	resp.Intent = strings.TrimSpace(strings.ToLower(resp.Intent))
	resp.Confidence = clampConfidence(resp.Confidence)
	if resp.Domain == "" || resp.Intent == "" {
		return nil, errors.New("missing fields")
	}

	return &resp, nil
}

func mapClassifiedProjectIntent(input string, resolved *ResolvedProject, intent string) chatbot.IntentMatch {
	match := chatbot.IntentMatch{Intent: chatbot.IntentUnknown, Confidence: 0.75}
	switch intent {
	case "project_create_plan":
		match.Intent = chatbot.IntentAddStagesAndTasks
	case "project_show_content", "project_show_plan", "project_list_tasks", "project_list_stages", "project_list_members", "project_count_done_tasks":
		match.Intent = chatbot.IntentShowProject
	}

	if resolved != nil && strings.TrimSpace(resolved.Title) != "" {
		match.ProjectTitle = resolved.Title
	}

	return match
}

func clampConfidence(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func containsProjectDecisionWords(lower string) bool {
	hasProjectWord := strings.Contains(lower, "project") || strings.Contains(lower, "projects") || strings.Contains(lower, "plan") || strings.Contains(lower, "stage") || strings.Contains(lower, "task")
	hasCreateWord := strings.Contains(lower, "create") || strings.Contains(lower, "add") || strings.Contains(lower, "start") || strings.Contains(lower, "new")

	if hasProjectWord && hasCreateWord {
		return true
	}

	if strings.Contains(lower, "project plan") {
		return true
	}

	return false
}

// isProjectListQuery detects natural-language questions about existing projects.
func isProjectListQuery(lower string) bool {
	if !(strings.Contains(lower, "project") || strings.Contains(lower, "projects")) {
		return false
	}

	triggers := []string{"what", "which", "list", "show", "do i have", "we have", "my", "any projects", "do we have", "do i have"}
	for _, t := range triggers {
		if strings.Contains(lower, t) {
			return true
		}
	}

	return false
}
