package chatbot

import (
	"regexp"
	"strings"

	"silencendo/services"
)

// Intent enumerations.
const (
	IntentCreateProject     = "create_project"
	IntentAddStagesAndTasks = "add_stages_tasks"
	IntentAssignResponsible = "assign_responsible"
	IntentListProjects      = "list_projects"
	IntentShowProject       = "show_project"
	IntentUnknown           = "unknown"
)

// IntentMatch captures parsed data from a message.
type IntentMatch struct {
	Intent       string
	ProjectTitle string
	Description  string
	StagePlans   []services.StagePlan
	EntityType   string
	EntityName   string
	AssigneeName string
	Filters      services.ListFilters
	Members      []string
	Confidence   float64
}

// DetectIntent performs lightweight rule-based intent detection and extraction.
func DetectIntent(message string) IntentMatch {
	msg := strings.TrimSpace(message)
	lower := strings.ToLower(msg)

	if isProjectListQuery(lower) {
		return IntentMatch{Intent: IntentListProjects, Confidence: 0.9}
	}

	// List projects. Project-related questions should be detected early and routed to DB.
	if strings.Contains(lower, "list projects") || strings.Contains(lower, "show projects") || strings.Contains(lower, "show my projects") || strings.Contains(lower, "projects list") || strings.Contains(lower, "what projects do i have") || strings.Contains(lower, "what projects do we have") || strings.Contains(lower, "what project do we have") || strings.Contains(lower, "my projects") || strings.Contains(lower, "/project list") || strings.Contains(lower, "project list") || strings.Contains(lower, "projects?") || strings.Contains(lower, "projects") || strings.Contains(lower, "список проектов") {
		return IntentMatch{Intent: IntentListProjects, Confidence: 0.9}
	}

	if isProjectQuestion(lower) {
		title := extractAfter(lower, []string{"show project", "project details", "покажи проект"})
		return IntentMatch{Intent: IntentShowProject, ProjectTitle: strings.TrimSpace(title), Confidence: 0.8}
	}

	if isPlanQuestion(lower) {
		return IntentMatch{Intent: IntentShowProject, Confidence: 0.85}
	}

	// Show project details.
	if strings.Contains(lower, "show project") || strings.Contains(lower, "project details") || strings.Contains(lower, "покажи проект") {
		title := extractAfter(lower, []string{"show project", "project details", "покажи проект"})
		return IntentMatch{Intent: IntentShowProject, ProjectTitle: strings.TrimSpace(title), Confidence: 0.7}
	}

	// Create project.
	if strings.Contains(lower, "create project") || strings.Contains(lower, "create a project") || strings.Contains(lower, "new project") || strings.Contains(lower, "создай проект") || strings.Contains(lower, "создать проект") {
		title := extractAfter(lower, []string{"create project", "create a project", "new project", "project", "создай проект", "создать проект"})
		return IntentMatch{Intent: IntentCreateProject, ProjectTitle: strings.TrimSpace(title), Confidence: 0.75}
	}

	// Add stages/tasks.
	if strings.Contains(lower, "create a plan") || strings.Contains(lower, "create plan") || strings.Contains(lower, "make a plan") || strings.Contains(lower, "step by step plan") || strings.Contains(lower, "сделай план") || strings.Contains(lower, "создай план") {
		return IntentMatch{Intent: IntentAddStagesAndTasks, ProjectTitle: extractProjectHint(lower), Members: extractPlanMembers(msg), Confidence: 0.7}
	}

	// Add stages/tasks.
	if strings.Contains(lower, "add stage") || strings.Contains(lower, "add stages") || strings.Contains(lower, "add task") || strings.Contains(lower, "add tasks") || strings.Contains(lower, "create stage") || strings.Contains(lower, "create task") {
		match := IntentMatch{Intent: IntentAddStagesAndTasks, Confidence: 0.7}
		match.StagePlans = parseStagesAndTasks(msg)
		match.ProjectTitle = extractProjectHint(lower)
		return match
	}

	// Assignment: if structured, assign; otherwise treat as planning for active project.
	if strings.Contains(lower, "assign") || strings.Contains(lower, "responsible") || strings.Contains(lower, "назначь") || containsAssignmentKeywords(lower) {
		entityType, assignee, entityName := parseAssign(lower)
		if entityType != "" && assignee != "" && entityName != "" {
			return IntentMatch{Intent: IntentAssignResponsible, EntityType: entityType, EntityName: entityName, AssigneeName: assignee, Confidence: 0.65}
		}
		return IntentMatch{Intent: IntentAddStagesAndTasks, ProjectTitle: extractProjectHint(lower), Members: extractPlanMembers(msg), Confidence: 0.6}
	}

	return IntentMatch{Intent: IntentUnknown, Confidence: 0.0}
}

func extractProjectHint(lower string) string {
	// naive extraction of project name after "project" keyword.
	projectRegex := regexp.MustCompile(`(?i)project\s+([^,;.]+)`)
	if m := projectRegex.FindStringSubmatch(lower); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	ruRegex := regexp.MustCompile(`(?i)проект\s+([^,;.]+)`)
	if m := ruRegex.FindStringSubmatch(lower); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func extractAfter(lower string, keywords []string) string {
	for _, kw := range keywords {
		if idx := strings.Index(lower, kw); idx >= 0 {
			after := lower[idx+len(kw):]
			after = strings.TrimSpace(after)
			after = strings.Trim(after, ":")
			return after
		}
	}
	return ""
}

// parseStagesAndTasks tries to recover simple stage/task lists like "stages: design; build | tasks: permit, plan".
func parseStagesAndTasks(input string) []services.StagePlan {
	var plans []services.StagePlan

	stageRegex := regexp.MustCompile(`(?i)stage[s]?:?([^|]+)`) // capture after stages:
	taskRegex := regexp.MustCompile(`(?i)task[s]?:?([^|]+)`)   // capture after tasks:

	stagePart := stageRegex.FindStringSubmatch(input)
	taskPart := taskRegex.FindStringSubmatch(input)

	if len(stagePart) == 2 {
		stageNames := splitItems(stagePart[1])
		for idx, name := range stageNames {
			plans = append(plans, services.StagePlan{Title: name, Order: idx + 1})
		}
	}

	// If tasks are present without explicit stages, attach to first stage if it exists.
	if len(taskPart) == 2 {
		tasks := splitItems(taskPart[1])
		var taskPlans []services.TaskPlan
		for _, t := range tasks {
			taskPlans = append(taskPlans, services.TaskPlan{Title: t})
		}
		if len(plans) == 0 {
			plans = append(plans, services.StagePlan{Title: "General", Order: 1, Tasks: taskPlans})
		} else {
			plans[0].Tasks = append(plans[0].Tasks, taskPlans...)
		}
	}

	return plans
}

func splitItems(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == '|' || r == '\n'
	})
	var out []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

// isProjectListQuery detects natural-language questions about existing projects.
func isProjectListQuery(lower string) bool {
	if !strings.Contains(lower, "projects") && !strings.Contains(lower, "project list") && !strings.Contains(lower, "/project list") && !strings.Contains(lower, "project do i have") && !strings.Contains(lower, "project do we have") && !strings.Contains(lower, "what project do we have") {
		return false
	}

	triggers := []string{"what", "which", "list", "show", "do i have", "we have", "my", "any projects", "do we have"}
	for _, t := range triggers {
		if strings.Contains(lower, t) {
			return true
		}
	}

	return false
}

func isProjectQuestion(lower string) bool {
	keywords := []string{
		"show project", "project details", "what tasks", "list stages", "members of project",
		"какие задачи", "этапы", "участники проекта", "покажи проект", "что в проекте", "какой статус проекта",
		"what content do we have", "content in this project", "project content", "content of this project", "what do we have in this project", "show the content", "what is the content", "content of",
	}
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func isPlanQuestion(lower string) bool {
	keywords := []string{"show the plan", "show plan", "what is the plan", "покажи план", "какой план", "view plan"}
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func extractPlanMembers(original string) []string {
	lower := strings.ToLower(original)
	idx := strings.Index(lower, "for ")
	if idx == -1 {
		idx = strings.Index(lower, "для ")
	}
	if idx == -1 {
		return nil
	}

	segment := original[idx+4:]
	// Stop at common delimiters like " to " or end of string.
	stopIdx := strings.Index(strings.ToLower(segment), " to ")
	if stopIdx != -1 {
		segment = segment[:stopIdx]
	}

	cleaned := strings.ReplaceAll(segment, " and ", ",")
	cleaned = strings.ReplaceAll(cleaned, " & ", ",")
	cleaned = strings.ReplaceAll(cleaned, " + ", ",")
	cleaned = strings.ReplaceAll(cleaned, " и ", ",")

	parts := strings.FieldsFunc(cleaned, func(r rune) bool {
		return r == ',' || r == ';' || r == '|' || r == '\n'
	})
	if len(parts) == 1 {
		// Fallback: split on whitespace when multiple tokens are packed without delimiters (e.g., "Omar Yussuf Fatima").
		words := strings.Fields(parts[0])
		if len(words) > 1 {
			parts = words
		}
	}
	var members []string
	for _, p := range parts {
		name := strings.TrimSpace(strings.TrimPrefix(p, "and"))
		name = strings.TrimSpace(strings.TrimPrefix(name, "и"))
		if name != "" {
			members = append(members, name)
		}
	}
	return members
}

func containsAssignmentKeywords(lower string) bool {
	// Check for common assignment expressions beyond basic keywords
	assignmentPatterns := []string{
		"give task", "give the task", "put on task", "put him on", "put her on",
		"have him do", "have her do", "let him", "let her", "make him",
		"make her", "should do", "needs to do", "take care of", "handle",
		"is responsible for", "will handle", "will take", "will work on",
	}
	for _, pattern := range assignmentPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

func parseAssign(lower string) (string, string, string) {
	// Examples: "assign John to stage Design", "assign Kate to task Permits".
	re := regexp.MustCompile(`assign\s+([^\s]+).*\s(stage|task)\s+([^,;.]+)`) // assignee, type, name
	if m := re.FindStringSubmatch(lower); len(m) == 4 {
		return m[2], m[1], strings.TrimSpace(m[3])
	}
	reRu := regexp.MustCompile(`назначь\s+([^\s]+).*\s(stage|task|этап|задачу)\s+([^,;.]+)`)
	if m := reRu.FindStringSubmatch(lower); len(m) == 4 {
		entityType := m[2]
		if entityType == "этап" {
			entityType = "stage"
		}
		if entityType == "задачу" {
			entityType = "task"
		}
		return entityType, m[1], strings.TrimSpace(m[3])
	}

	// More flexible assignment patterns for natural language
	// Handle "Yusuf to craft eyes of ender" (where "assign Yusuf to craft eyes of ender" was intended)
	flexRe := regexp.MustCompile(`(?:assign|put|give|let|make|have|get)\s+(?:to\s+)?(\w+)\s+(?:on\s+|with\s+|to\s+|for\s+)?(?:task|the\s+task|stage|the\s+stage)?\s*(.+?)(?:\.|$|,|!)`)
	if m := flexRe.FindStringSubmatch(lower); len(m) >= 3 {
		// Try to determine if it's a task or stage based on context clues
		entityName := strings.TrimSpace(m[2])
		entityType := "task" // default to task

		// Check if the entity name contains stage-related terms
		if strings.Contains(entityName, "stage") || strings.Contains(entityName, "phase") || strings.Contains(entityName, "step") {
			entityType = "stage"
		}

		return entityType, m[1], entityName
	}

	// Handle patterns like "Yusuf should do craft eyes of ender", "Yusuf will handle the task", etc.
	verbBasedRe := regexp.MustCompile(`(\w+)\s+(?:should|will|needs to|has to|must|can|could|would|is responsible for|will take on|will work on|will complete|will finish|take care of|handles?)\s+(?:the\s+)?(task|stage)?\s*(.+?)(?:\.|$|,|!)`)
	if m := verbBasedRe.FindStringSubmatch(lower); len(m) >= 3 {
		entityType := "task" // default to task
		entityName := strings.TrimSpace(m[3])

		if len(m) > 2 && strings.TrimSpace(m[2]) != "" {
			entityType = strings.TrimSpace(m[2])
		} else {
			// Determine type based on content
			if strings.Contains(entityName, "stage") || strings.Contains(entityName, "phase") || strings.Contains(entityName, "step") {
				entityType = "stage"
			}
		}

		return entityType, m[1], entityName
	}

	// Handle patterns like "make Dastan and Yussuf responsible of setting spawn and placing beds"
	multiAssignRe := regexp.MustCompile(`(?:make|let|put|assign)\s+(.+?)\s+(?:responsible|in charge|to handle|to do|for)\s+(?:of\s+)?(.+?)(?:\.|$|,|!)`)
	if m := multiAssignRe.FindStringSubmatch(lower); len(m) >= 3 {
		// For multiple assignees, we'll take the first one for now
		// Split by "and" or "," to get the first assignee
		assignees := strings.Split(m[1], " and ")
		firstAssignee := strings.TrimSpace(assignees[0])
		// Remove any commas
		commaSplit := strings.Split(firstAssignee, ",")
		firstAssignee = strings.TrimSpace(commaSplit[0])

		entityName := strings.TrimSpace(m[2])
		entityType := "task" // default to task

		// Determine type based on content
		if strings.Contains(entityName, "stage") || strings.Contains(entityName, "phase") || strings.Contains(entityName, "step") {
			entityType = "stage"
		}

		return entityType, firstAssignee, entityName
	}

	return "", "", ""
}
