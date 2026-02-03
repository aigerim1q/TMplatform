package chatbot

import (
	"context"
	"regexp"
	"strings"

	router "silencendo/intent_router"
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
	TaskID       int
	Assignees    []string
}

// Global router instance - needs to be initialized externally
var globalRouter *router.Router

// SetGlobalRouter sets the global LLM router to be used for intent detection
func SetGlobalRouter(r *router.Router) {
	globalRouter = r
}

// DetectIntent uses the LLM as the central reasoning engine for intent detection and extraction.
// If no LLM router is set, falls back to rule-based detection.
func DetectIntent(message string) IntentMatch {
	if globalRouter != nil {
		// Use LLM-based intent detection
		ctx := context.Background()
		result := DetectIntentWithLLM(ctx, globalRouter, message)

		// If the LLM detected an assignment intent but with incomplete information,
		// fall back to rule-based parsing to supplement the missing data
		if result.Intent == IntentAssignResponsible {
			entityName := strings.TrimSpace(result.EntityName)
			assigneeName := strings.TrimSpace(result.AssigneeName)

			// If either entity name or assignee name is missing, try rule-based parsing
			if entityName == "" || assigneeName == "" {
				ruleResult := detectIntentRuleBased(message)

				// Only use rule-based result if it provides more complete information
				if ruleResult.Intent == IntentAssignResponsible {
					// If rule-based parsing has both pieces of information that LLM missed, use it
					if entityName == "" && ruleResult.EntityName != "" {
						result.EntityName = ruleResult.EntityName
					}
					if assigneeName == "" && ruleResult.AssigneeName != "" {
						result.AssigneeName = ruleResult.AssigneeName
					}

					// If the rule-based parsing provides better assignee information
					// (e.g., catches multiple assignees separated by commas), prioritize it
					if ruleResult.AssigneeName != "" {
						result.AssigneeName = ruleResult.AssigneeName
					}

					// Update confidence if we're supplementing with rule-based info
					if result.EntityName != "" && result.AssigneeName != "" {
						result.Confidence = max(result.Confidence, ruleResult.Confidence)
					}
				}
			}

		}

		// Additional check: if the LLM didn't detect assignment but the message clearly looks like an assignment,
		// use rule-based parsing
		if result.Intent != IntentAssignResponsible {
			// Check if the message contains clear assignment patterns
			lowerMsg := strings.ToLower(message)
			if strings.Contains(lowerMsg, "assign") && (strings.Contains(lowerMsg, "#") || strings.Contains(lowerMsg, "task") || strings.Contains(lowerMsg, "stage")) {
				ruleResult := detectIntentRuleBased(message)
				if ruleResult.Intent == IntentAssignResponsible && ruleResult.EntityName != "" && ruleResult.AssigneeName != "" {
					// Use the rule-based result if it successfully parsed an assignment
					return ruleResult
				}
			}
		}

		return result
	}

	// Fallback to rule-based detection if no LLM router is available
	return detectIntentRuleBased(message)
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// detectIntentRuleBased is the original rule-based implementation for fallback purposes
func detectIntentRuleBased(message string) IntentMatch {
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

	// Check for stable ID patterns in assignment context
	stableIDPattern := regexp.MustCompile(`#\s*\d+`)
	hasStableID := stableIDPattern.MatchString(lower)

	// Assignment: if structured, assign; otherwise treat as planning for active project.
	if strings.Contains(lower, "assign") || strings.Contains(lower, "responsible") || strings.Contains(lower, "назначь") || containsAssignmentKeywords(lower) || hasStableID {
		entityType, assignee, entityName := parseAssign(lower)

		// If we have a stable ID in entityName, accept it as a valid assignment even if traditional parsing fails
		if hasStableID && entityName != "" && assignee != "" {
			if entityType == "" {
				entityType = "task" // Default to task when using stable IDs
			}
			return IntentMatch{Intent: IntentAssignResponsible, EntityType: entityType, EntityName: entityName, AssigneeName: assignee, Confidence: 0.75}
		}

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

var (
	assignTaskRegex = regexp.MustCompile(`(?i)assign\s+(.+?)\s+(?:to\s+)?(?:task\s+)?#?([\w-]+)`)
)

func parseAssign(input string) (string, string, string) {
	inputLower := strings.ToLower(input)

	// Handle common typos in "assign" by replacing them
	// Replace "assing" with "assign" and other potential typos
	inputLower = strings.ReplaceAll(inputLower, " assing ", " assign ")
	if strings.HasPrefix(inputLower, "assing") {
		inputLower = strings.Replace(inputLower, "assing", "assign", 1)
	}

	// Normalize input to handle variations
	normalized := inputLower

	// Pattern to match: assign <names> to task #<id> or assign <names> to task <id>
	// This includes variations with commas between names
	assignPattern := regexp.MustCompile(`(?:assign|assing|make|let|put)\s+(.+?)\s+to\s+(?:the\s+)?(?:task|stage)?\s*#?(\d+)`)
	if m := assignPattern.FindStringSubmatch(normalized); len(m) >= 3 {
		assigneesStr := m[1]
		taskID := m[2]

		// Process assignees - split by "and" and "," and clean up
		assigneesStr = strings.ReplaceAll(assigneesStr, " and ", ",")
		assigneesStr = strings.ReplaceAll(assigneesStr, " & ", ",")
		assigneesStr = strings.ReplaceAll(assigneesStr, " + ", ",")

		assigneeParts := strings.FieldsFunc(assigneesStr, func(r rune) bool {
			return r == ','
		})

		// Clean up individual assignees
		var cleanAssignees []string
		for _, part := range assigneeParts {
			cleanPart := strings.TrimSpace(part)
			if cleanPart != "" && cleanPart != "to" && cleanPart != "the" && cleanPart != "a" && cleanPart != "an" {
				// Remove trailing punctuation
				cleanPart = strings.Trim(cleanPart, " .,!?;")
				if cleanPart != "" {
					cleanAssignees = append(cleanAssignees, cleanPart)
				}
			}
		}

		if len(cleanAssignees) > 0 {
			// Join multiple assignees with " and "
			assigneeName := strings.Join(cleanAssignees, " and ")
			entityName := "#" + taskID // Use the stable ID format
			return "task", assigneeName, entityName
		}
	}

	// Additional pattern to catch "assign X, Y to task #Z" format specifically
	commaAssignPattern := regexp.MustCompile(`(?:assign|assing|make|let|put)\s+([^,]+(?:,\s*[^,]+)+)\s+to\s+(?:the\s+)?(?:task|stage)?\s*#?(\d+)`)
	if m := commaAssignPattern.FindStringSubmatch(normalized); len(m) >= 3 {
		assigneesStr := m[1]
		taskID := m[2]

		// Process assignees - split by commas and clean up
		assigneeParts := strings.FieldsFunc(assigneesStr, func(r rune) bool {
			return r == ','
		})

		// Clean up individual assignees
		var cleanAssignees []string
		for _, part := range assigneeParts {
			cleanPart := strings.TrimSpace(part)
			if cleanPart != "" && cleanPart != "to" && cleanPart != "the" && cleanPart != "a" && cleanPart != "an" {
				// Remove trailing punctuation
				cleanPart = strings.Trim(cleanPart, " .,!?;")
				if cleanPart != "" {
					cleanAssignees = append(cleanAssignees, cleanPart)
				}
			}
		}

		if len(cleanAssignees) > 0 {
			// Join multiple assignees with " and "
			assigneeName := strings.Join(cleanAssignees, " and ")
			entityName := "#" + taskID // Use the stable ID format
			return "task", assigneeName, entityName
		}
	}

	// Additional pattern to catch "assign X, Y to #Z" format (without "task" word)
	simpleCommaAssignPattern := regexp.MustCompile(`(?:assign|assing|make|let|put)\s+([^,]+(?:,\s*[^,]+)*,\s*[^,]+|\w+(?:\s+and\s+\w+|\s*&\s*\w+|\s+\+\s+\w+)*)\s+to\s+(?:task\s+)?#?(\d+)`)
	if m := simpleCommaAssignPattern.FindStringSubmatch(normalized); len(m) >= 3 {
		assigneesStr := m[1]
		taskID := m[2]

		// Process assignees - split by commas and clean up
		assigneeParts := strings.FieldsFunc(assigneesStr, func(r rune) bool {
			return r == ','
		})

		// Clean up individual assignees
		var cleanAssignees []string
		for _, part := range assigneeParts {
			cleanPart := strings.TrimSpace(part)
			if cleanPart != "" && cleanPart != "to" && cleanPart != "the" && cleanPart != "a" && cleanPart != "an" {
				// Remove trailing punctuation
				cleanPart = strings.Trim(cleanPart, " .,!?;")
				if cleanPart != "" {
					cleanAssignees = append(cleanAssignees, cleanPart)
				}
			}
		}

		// If no commas were used, check for "and", "&", or "+" separators
		if len(cleanAssignees) == 0 {
			assigneesStr = strings.ReplaceAll(assigneesStr, " and ", ",")
			assigneesStr = strings.ReplaceAll(assigneesStr, " & ", ",")
			assigneesStr = strings.ReplaceAll(assigneesStr, " + ", ",")

			assigneeParts = strings.FieldsFunc(assigneesStr, func(r rune) bool {
				return r == ','
			})

			for _, part := range assigneeParts {
				cleanPart := strings.TrimSpace(part)
				if cleanPart != "" && cleanPart != "to" && cleanPart != "the" && cleanPart != "a" && cleanPart != "an" {
					// Remove trailing punctuation
					cleanPart = strings.Trim(cleanPart, " .,!?;")
					if cleanPart != "" {
						cleanAssignees = append(cleanAssignees, cleanPart)
					}
				}
			}
		}

		if len(cleanAssignees) > 0 {
			// Join multiple assignees with " and "
			assigneeName := strings.Join(cleanAssignees, " and ")
			entityName := "#" + taskID // Use the stable ID format
			return "task", assigneeName, entityName
		}
	}

	// Detect stable task ID format like "#8", "# 8", etc.
	stableIDRegex := regexp.MustCompile(`#\s*(\d+)`)
	stableIDMatches := stableIDRegex.FindStringSubmatch(input)

	if len(stableIDMatches) >= 2 {
		stableID := stableIDMatches[1]
		taskRef := "#" + stableID // Construct the reference like "#8"

		// Extract assignees from before the "to" part or from the whole string
		toIndex := strings.Index(input, " to ")
		var assigneesStr string

		if toIndex != -1 {
			// Get the part before "to" which should contain assignees
			beforeTo := input[:toIndex]
			// Remove the stable ID reference from the beforeTo part to avoid confusion
			cleanBeforeTo := strings.TrimSpace(strings.ReplaceAll(beforeTo, taskRef, ""))

			// Extract assignees by looking for patterns like "assign X and Y" or just "X and Y"
			// Look for common assignment verbs and extract what comes after
			assignPattern := regexp.MustCompile(`(?:assign|make|let|put|get|have|add)\s+(.+)`)
			assignMatches := assignPattern.FindStringSubmatch(cleanBeforeTo)

			if len(assignMatches) >= 2 {
				assigneesStr = assignMatches[1]
			} else {
				// If no verb pattern found, just use the cleaned beforeTo part
				assigneesStr = cleanBeforeTo
			}
		} else {
			// If no "to" found, extract assignees from the whole string except the stable ID
			cleanLower := strings.TrimSpace(strings.ReplaceAll(input, taskRef, ""))
			// Look for assignment patterns
			assignPattern := regexp.MustCompile(`(?:assign|make|let|put|get|have|add)\s+(.+?)\s+#\s*\d+`)
			assignMatches := assignPattern.FindStringSubmatch(input)

			if len(assignMatches) >= 2 {
				assigneesStr = assignMatches[1]
			} else {
				// As a fallback, just take everything that's not the stable ID and common words
				words := strings.Fields(cleanLower)
				var extractedAssignees []string
				for _, word := range words {
					trimmed := strings.Trim(word, ",.!;()")
					// Skip common verbs and keywords
					if trimmed == "assign" || trimmed == "make" || trimmed == "let" || trimmed == "put" ||
						trimmed == "get" || trimmed == "have" || trimmed == "add" || trimmed == "to" ||
						trimmed == "task" || trimmed == "taks" || trimmed == "the" {
						continue
					}
					// Skip if it's a number or looks like a stable ID
					if regexp.MustCompile(`^\d+$`).MatchString(trimmed) || strings.Contains(trimmed, "#") {
						continue
					}
					extractedAssignees = append(extractedAssignees, trimmed)
				}
				assigneesStr = strings.Join(extractedAssignees, " ")
			}
		}

		// Process assignees - split by "and" and "," and clean up
		assigneesStr = strings.ReplaceAll(assigneesStr, " and ", ",")
		assigneesStr = strings.ReplaceAll(assigneesStr, " & ", ",")
		assigneesStr = strings.ReplaceAll(assigneesStr, " + ", ",")

		assigneeParts := strings.FieldsFunc(assigneesStr, func(r rune) bool {
			return r == ','
		})

		// Clean up individual assignees
		var cleanAssignees []string
		for _, part := range assigneeParts {
			cleanPart := strings.TrimSpace(part)
			if cleanPart != "" && cleanPart != "to" && cleanPart != "" {
				cleanAssignees = append(cleanAssignees, cleanPart)
			}
		}

		// Return the first assignee for compatibility with current interface
		if len(cleanAssignees) > 0 {
			return "task", cleanAssignees[0], taskRef // Use the stable ID reference
		}

		// If we couldn't parse assignees but have a stable ID, return what we have
		return "task", "", taskRef
	}

	// Examples: "assign John to stage Design", "assign Kate to task Permits".
	re := regexp.MustCompile(`assign\s+([^\s]+).*\s(stage|task)\s+([^,;.]+)`) // assignee, type, name
	if m := re.FindStringSubmatch(inputLower); len(m) == 4 {
		return m[2], m[1], strings.TrimSpace(m[3])
	}

	// Handle "assign X and Y to task Z" - multiple assignees
	multiAssignRe := regexp.MustCompile(`assign\s+(.+?)\s+to\s+(?:the\s+)?(task|stage)?\s*(.+?)(?:\.|$|,|!)`)
	if m := multiAssignRe.FindStringSubmatch(inputLower); len(m) >= 3 {
		// Extract first assignee from "X and Y"
		assignees := strings.Split(m[1], " and ")
		firstAssignee := strings.TrimSpace(assignees[0])

		entityName := strings.TrimSpace(m[3])
		entityType := "task" // default to task

		if len(m) > 2 && strings.TrimSpace(m[2]) != "" {
			entityType = strings.TrimSpace(m[2])
		} else {
			// Determine type based on content
			if strings.Contains(entityName, "stage") || strings.Contains(entityName, "phase") || strings.Contains(entityName, "step") {
				entityType = "stage"
			}
		}

		return entityType, firstAssignee, entityName
	}

	reRu := regexp.MustCompile(`назначь\s+([^\s]+).*\s(stage|task|этап|задачу)\s+([^,;.]+)`)
	if m := reRu.FindStringSubmatch(inputLower); len(m) == 4 {
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
	if m := flexRe.FindStringSubmatch(inputLower); len(m) >= 3 {
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
	if m := verbBasedRe.FindStringSubmatch(inputLower); len(m) >= 3 {
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
	multiAssignReOld := regexp.MustCompile(`(?:make|let|put|assign)\s+(.+?)\s+(?:responsible|in charge|to handle|to do|for)\s+(?:of\s+)?(.+?)(?:\.|$|,|!)`)
	if m := multiAssignReOld.FindStringSubmatch(inputLower); len(m) >= 3 {
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

	// Handle partial assignment commands like "assign Yusuf and Dastan" (without task specification)
	partialAssignRe := regexp.MustCompile(`(?:assign|make|let|put|get|have|add)\s+(.+?)(?:\.|$|,|!)`)
	if m := partialAssignRe.FindStringSubmatch(inputLower); len(m) >= 2 {
		// Check that this doesn't also contain "to" followed by a task, to avoid conflicting with other patterns
		if !strings.Contains(inputLower, " to ") {
			assigneesStr := strings.TrimSpace(m[1])
			// Process assignees - split by "and" and "," and clean up
			assigneesStr = strings.ReplaceAll(assigneesStr, " and ", ",")
			assigneesStr = strings.ReplaceAll(assigneesStr, " & ", ",")
			assigneesStr = strings.ReplaceAll(assigneesStr, " + ", ",")

			assigneeParts := strings.FieldsFunc(assigneesStr, func(r rune) bool {
				return r == ','
			})

			// Clean up individual assignees
			var cleanAssignees []string
			for _, part := range assigneeParts {
				cleanPart := strings.TrimSpace(part)
				if cleanPart != "" && cleanPart != "to" && cleanPart != "" {
					cleanPart = strings.Trim(cleanPart, " .,!?;")
					if cleanPart != "" {
						cleanAssignees = append(cleanAssignees, cleanPart)
					}
				}
			}

			// Return the first assignee if we found any
			if len(cleanAssignees) > 0 {
				return "task", cleanAssignees[0], "" // Empty entityName indicates missing task
			}
		}
	}

	return "", "", ""
}
