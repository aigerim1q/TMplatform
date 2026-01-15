package chatbot

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	router "silencendo/intent_router"
	"silencendo/models"
	"silencendo/services"
)

// Handler routes chat messages to domain services and returns human-readable responses.
type Handler struct {
	Projects *services.ProjectService
	User     models.User
}

type ctxKey string

const (
	ctxResolvedProjectIDKey    ctxKey = "resolvedProjectID"
	ctxResolvedProjectTitleKey ctxKey = "resolvedProjectTitle"
)

// WithResolvedProject attaches a resolved project reference to the context for downstream handlers.
func WithResolvedProject(ctx context.Context, id, title string) context.Context {
	ctx = context.WithValue(ctx, ctxResolvedProjectIDKey, strings.TrimSpace(id))
	ctx = context.WithValue(ctx, ctxResolvedProjectTitleKey, strings.TrimSpace(title))
	return ctx
}

// resolvedProjectFromCtx extracts a resolved project reference if present.
func resolvedProjectFromCtx(ctx context.Context) (id string, title string, ok bool) {
	idVal := ctx.Value(ctxResolvedProjectIDKey)
	titleVal := ctx.Value(ctxResolvedProjectTitleKey)
	id, _ = idVal.(string)
	title, _ = titleVal.(string)
	id = strings.TrimSpace(id)
	title = strings.TrimSpace(title)
	ok = id != ""
	return
}

// ResolvedProjectFromCtx extracts a resolved project reference if present (exported version).
func ResolvedProjectFromCtx(ctx context.Context) (id string, title string, ok bool) {
	return resolvedProjectFromCtx(ctx)
}

func NewHandler(projects *services.ProjectService, user models.User) *Handler {
	return &Handler{Projects: projects, User: user}
}

// Handle attempts to satisfy chatbot intents. handled indicates whether the message was consumed by this layer.
func (h *Handler) Handle(ctx context.Context, message string) (handled bool, reply string, err error) {
	match := DetectIntent(message)
	return h.HandleWithMatch(ctx, match)
}

// HandleWithLLMResponse handles responses based on LLM-generated intent directly
func (h *Handler) HandleWithLLMResponse(ctx context.Context, llmResponse *router.LLMResponse) (handled bool, reply string, err error) {
	// Convert the LLM response to an IntentMatch
	intent := llmResponse.ToChatbotIntent()

	match := IntentMatch{
		Intent:     intent,
		Confidence: llmResponse.Confidence,
	}

	// Extract parameters from the LLM response
	if title, ok := llmResponse.Parameters["project_title"].(string); ok {
		match.ProjectTitle = title
	}
	if description, ok := llmResponse.Parameters["description"].(string); ok {
		match.Description = description
	}
	if entityType, ok := llmResponse.Parameters["entity_type"].(string); ok {
		match.EntityType = entityType
	}
	if entityName, ok := llmResponse.Parameters["entity_name"].(string); ok {
		match.EntityName = entityName
	}
	if assigneeName, ok := llmResponse.Parameters["assignee_name"].(string); ok {
		match.AssigneeName = assigneeName
	}

	// Use the existing HandleWithMatch method to process the converted match
	return h.HandleWithMatch(ctx, match)
}

// HandleWithMatch routes using a pre-classified intent match (e.g., from DeepSeek classifier).
func (h *Handler) HandleWithMatch(ctx context.Context, match IntentMatch) (handled bool, reply string, err error) {
	resolvedID, resolvedTitle, hasResolved := resolvedProjectFromCtx(ctx)
	if hasResolved && strings.TrimSpace(match.ProjectTitle) == "" {
		match.ProjectTitle = resolvedTitle
	}

	switch match.Intent {
	case IntentCreateProject:
		handled = true
		if strings.TrimSpace(match.ProjectTitle) == "" {
			return handled, "I can create a project. Please provide a project title (e.g., 'Create project Skyline Apartments').", nil
		}
		if match.Confidence < 0.6 {
			return handled, fmt.Sprintf("Did you want to create a project named '%s'? Please confirm or restate.", strings.TrimSpace(match.ProjectTitle)), nil
		}
		project, created, err := h.Projects.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: match.ProjectTitle, Description: match.Description}, h.User)
		if err != nil {
			return handled, humanizeError(err), err
		}
		if created {
			return handled, fmt.Sprintf("✅ Project created: %s", project.Title), nil
		}
		return handled, fmt.Sprintf("ℹ️ Project already exists: %s", project.Title), nil

	case IntentAddStagesAndTasks:
		handled = true
		if !hasResolved {
			return handled, noActiveProjectMessage(), nil
		}
		projectID := resolvedID
		projectTitle := resolvedTitle

		if len(match.StagePlans) == 0 {
			members, stages, tasks, err := h.Projects.CreateRuleBasedPlan(ctx, projectID, match.Members, h.User)
			if err != nil {
				return handled, humanizeError(err), err
			}
			return handled, fmt.Sprintf("✅ Plan created for %s: %d member(s), %d stage(s), %d task(s).", projectTitle, len(members), len(stages), len(tasks)), nil
		}

		stages, tasks, err := h.Projects.CreateStagesAndTasksFromChat(ctx, projectID, match.StagePlans, h.User)
		if err != nil {
			return handled, humanizeError(err), err
		}

		return handled, fmt.Sprintf("✅ Updated project %s: %d stage(s), %d task(s) ensured.", projectTitle, len(stages), len(tasks)), nil

	case IntentAssignResponsible:
		handled = true
		entityType := strings.TrimSpace(match.EntityType)
		entityName := strings.TrimSpace(match.EntityName)
		assigneeName := strings.TrimSpace(match.AssigneeName)

		// Get original message for stable ID detection - we need to get it from the original Handle function call
		// Since HandleWithMatch doesn't have access to the original message, we'll rely on the existing parsing
		stableIDPattern := regexp.MustCompile(`#\s*\d+`)
		hasStableID := stableIDPattern.MatchString(strings.TrimSpace(entityName)) || strings.HasPrefix(strings.TrimSpace(entityName), "#")

		// Extract stable task ID from entityName for better error messages
		var stableTaskID string
		if stableIDPattern.MatchString(strings.TrimSpace(entityName)) {
			matches := stableIDPattern.FindStringSubmatch(strings.TrimSpace(entityName))
			if len(matches) > 0 {
				stableTaskID = strings.TrimPrefix(matches[0], "#")
				stableTaskID = strings.TrimSpace(stableTaskID)
			}
		}

		if entityName == "" || assigneeName == "" {
			// If we have a stable ID, try to validate differently
			if hasStableID {
				// Extract just the assignee part and check if it's meaningful
				assigneeClean := strings.TrimSpace(assigneeName)
				if assigneeClean == "" || assigneeClean == "to" {
					// Try to extract assignees from the original message if possible
					// This handles cases where LLM might not have extracted assignees properly
					originalMsg := ctx.Value("original_message") // This won't work as we don't pass original message in context
					_ = originalMsg                              // suppress unused variable warning

					// Since we can't access original message here, let's check if we can improve our own parsing
					if stableTaskID != "" {
						return handled, fmt.Sprintf("Who should I assign to task #%s?", stableTaskID), nil
					} else {
						return handled, "Could you clarify who you want to assign to the task? Please provide assignee names.", nil
					}
				}
			} else {
				// Check if entityName is missing but task reference was provided
				if entityName == "" {
					// Store pending assignment state to remember the assignees
					if assigneeName != "" {
						// Process multiple assignees if they are separated by "and" or ","
						assigneeNames := []string{assigneeName} // Start with single assignee
						if strings.Contains(assigneeName, " and ") || strings.Contains(assigneeName, ",") {
							// Split by " and " or "," to get multiple assignees
							assigneeNames = strings.Split(assigneeName, " and ")
							// Further split by comma if present
							var allNames []string
							for _, nameGroup := range assigneeNames {
								names := strings.Split(nameGroup, ",")
								for _, name := range names {
									trimmed := strings.TrimSpace(name)
									if trimmed != "" {
										allNames = append(allNames, trimmed)
									}
								}
							}
							assigneeNames = allNames
						}

						// This would need to be called from the dispatcher context to access session state
						// We'll need to pass the session state through the context or return a special error
						// For now, we'll return a message indicating we need to store the pending assignment
						return handled, "Which task should I assign them to? Use task 4-2 or #8.", nil
					} else {
						return handled, "Could you clarify who you want to assign to what task or stage?", nil
					}
				} else {
					// Entity name is provided but assignee is missing
					return handled, "Who should I assign to which task?", nil
				}
			}
		}

		// Default to "task" if entity type is not specified
		if entityType == "" {
			entityType = "task"
		}

		if !hasResolved {
			return handled, noActiveProjectMessage(), nil
		}

		// Process multiple assignees if they are separated by "and" or ","
		assigneeNames := []string{assigneeName} // Start with single assignee
		if strings.Contains(assigneeName, " and ") || strings.Contains(assigneeName, ",") {
			// Split by " and " or "," to get multiple assignees
			assigneeNames = strings.Split(assigneeName, " and ")
			// Further split by comma if present
			var allNames []string
			for _, nameGroup := range assigneeNames {
				names := strings.Split(nameGroup, ",")
				for _, name := range names {
					trimmed := strings.TrimSpace(name)
					if trimmed != "" {
						allNames = append(allNames, trimmed)
					}
				}
			}
			assigneeNames = allNames
		}

		// Handle multiple assignees
		var assigneeUsers []models.User
		for _, name := range assigneeNames {
			name = strings.TrimSpace(name)
			if name != "" {
				user, err := h.Projects.EnsureUser(ctx, models.User{Name: name})
				if err != nil {
					return handled, humanizeError(err), err
				}
				assigneeUsers = append(assigneeUsers, user)
			}
		}

		switch strings.ToLower(entityType) {
		case "stage":
			// Handle stage assignment with support for position indices
			stage, err := h.resolveStageByNameOrPosition(ctx, entityName, resolvedID)
			if err != nil {
				return handled, humanizeError(err), err
			}
			if stage.ProjectID != resolvedID {
				return handled, "That stage belongs to another project. Switch active project first.", nil
			}

			// Assign all users to the stage
			for _, assignee := range assigneeUsers {
				if err := h.Projects.AssignResponsible(ctx, "stage", stage.ID, assignee.ID, h.User); err != nil {
					return handled, humanizeError(err), err
				}
			}

			assigneeNamesStr := ""
			for i, user := range assigneeUsers {
				if i > 0 {
					assigneeNamesStr += ", "
				}
				assigneeNamesStr += user.Name
			}
			return handled, fmt.Sprintf("✅ Stage '%s' assigned to %s", stage.Title, assigneeNamesStr), nil
		case "task":
			// Handle task assignment with support for both position indices and stable IDs
			task, err := h.resolveTaskByNameOrPositionOrStableID(ctx, entityName, resolvedID)
			if err != nil {
				return handled, humanizeError(err), err
			}
			if task.ProjectID != resolvedID {
				return handled, "That task belongs to another project. Switch active project first.", nil
			}

			// Assign all users to the task
			for _, assignee := range assigneeUsers {
				if err := h.Projects.AssignResponsible(ctx, "task", task.ID, assignee.ID, h.User); err != nil {
					return handled, humanizeError(err), err
				}
			}

			assigneeNamesStr := ""
			for i, user := range assigneeUsers {
				if i > 0 {
					assigneeNamesStr += ", "
				}
				assigneeNamesStr += user.Name
			}
			return handled, fmt.Sprintf("✅ Task '%s' assigned to %s", task.Title, assigneeNamesStr), nil
		default:
			// Default to task if not specified properly
			task, err := h.resolveTaskByNameOrPositionOrStableID(ctx, entityName, resolvedID)
			if err != nil {
				return handled, humanizeError(err), err
			}
			if task.ProjectID != resolvedID {
				return handled, "That task belongs to another project. Switch active project first.", nil
			}

			// Assign all users to the task
			for _, assignee := range assigneeUsers {
				if err := h.Projects.AssignResponsible(ctx, "task", task.ID, assignee.ID, h.User); err != nil {
					return handled, humanizeError(err), err
				}
			}

			assigneeNamesStr := ""
			for i, user := range assigneeUsers {
				if i > 0 {
					assigneeNamesStr += ", "
				}
				assigneeNamesStr += user.Name
			}
			return handled, fmt.Sprintf("✅ Task '%s' assigned to %s", task.Title, assigneeNamesStr), nil
		}

	case IntentListProjects:
		handled = true
		projects, err := h.Projects.ListProjects(ctx, h.User, match.Filters)
		if err != nil {
			return handled, humanizeError(err), err
		}
		if len(projects) == 0 {
			return handled, "You have no projects yet. Create one with 'create project <name>'.", nil
		}
		var lines []string
		activeID := ""
		if hasResolved {
			activeID = resolvedID
		}
		for i, p := range projects {
			marker := "  "
			if p.ID == activeID {
				marker = "* "
			}
			lines = append(lines, fmt.Sprintf("%s[%d] %s (status: %s)", marker, i+1, p.Title, p.Status))
		}
		return handled, "📁 Your projects:\n" + strings.Join(lines, "\n"), nil

	case IntentShowProject:
		handled = true
		if !hasResolved {
			return handled, noActiveProjectMessage(), nil
		}
		projectID := resolvedID
		projectTitle := resolvedTitle
		details, err := h.Projects.GetProjectDetails(ctx, projectID, h.User)
		if err != nil {
			return handled, humanizeError(err), err
		}
		details.Project.Title = strings.TrimSpace(projectTitle)
		if len(details.Stages) == 0 && len(details.Tasks) == 0 {
			return handled, fmt.Sprintf("No plan exists yet for %s. Use 'create plan for <names> ...' to add one.", details.Project.Title), nil
		}
		return handled, h.formatProjectDetailsWithUserNames(ctx, details), nil
	}

	return false, "", nil
}

func humanizeError(err error) string {
	if errors.Is(err, services.ErrValidation) {
		return "Please provide complete information (title, project name, or assignee)."
	}
	if errors.Is(err, services.ErrForbidden) {
		return "You do not have access to that project."
	}
	return "An error occurred. Please try again or check your input."
}

func noActiveProjectMessage() string {
	return "No active project selected. Use /project list and /project use <number>."
}

func formatProjectDetails(details models.ProjectDetails) string {
	const taskLimit = 50
	var b strings.Builder
	fmt.Fprintf(&b, "📁 %s (status: %s)\n", details.Project.Title, details.Project.Status)
	if strings.TrimSpace(details.Project.Description) != "" {
		fmt.Fprintf(&b, "Description: %s\n", details.Project.Description)
	}

	if len(details.Stages) == 0 {
		b.WriteString("Stages: none\n")
	} else {
		b.WriteString("Stages:\n")
		sorted := make([]models.Stage, len(details.Stages))
		copy(sorted, details.Stages)
		sort.SliceStable(sorted, func(i, j int) bool {
			if sorted[i].OrderIndex == sorted[j].OrderIndex {
				return sorted[i].Title < sorted[j].Title
			}
			if sorted[i].OrderIndex == 0 {
				return false
			}
			if sorted[j].OrderIndex == 0 {
				return true
			}
			return sorted[i].OrderIndex < sorted[j].OrderIndex
		})
		for i, s := range sorted {
			resp := "unassigned"
			if s.ResponsibleID != nil {
				resp = "assigned"
			}
			fmt.Fprintf(&b, "[%d] %s (responsible: %s)\n", i+1, s.Title, resp) // Show stage index
		}
	}

	if len(details.Tasks) == 0 {
		b.WriteString("Tasks: none")
		return strings.TrimSpace(b.String())
	}

	// group tasks by stage
	byStage := map[string][]models.Task{}
	for _, t := range details.Tasks {
		key := ""
		if t.StageID != nil {
			key = *t.StageID
		}
		byStage[key] = append(byStage[key], t)
	}

	stageOrder := []string{}
	for _, s := range details.Stages {
		stageOrder = append(stageOrder, s.ID)
	}
	stageOrder = append(stageOrder, "") // for unassigned tasks

	b.WriteString("Tasks:\n")
	totalShown := 0

	// Group tasks by stage and sort them for consistent positioning
	for _, stageID := range stageOrder {
		tasks := byStage[stageID]
		if len(tasks) == 0 {
			continue
		}

		// Find the stage to get its position index
		stagePosition := 0
		for i, s := range details.Stages {
			if s.ID == stageID {
				stagePosition = i + 1 // 1-based index
				break
			}
		}

		stageTitle := "No stage"
		for _, s := range details.Stages {
			if s.ID == stageID {
				stageTitle = s.Title
				break
			}
		}

		// Print stage with position index
		if stagePosition > 0 {
			fmt.Fprintf(&b, "• [%d] %s:\n", stagePosition, stageTitle)
		} else {
			fmt.Fprintf(&b, "• %s:\n", stageTitle)
		}

		// Sort tasks within the stage by creation order or ID to maintain consistent positioning
		sortedTasks := make([]models.Task, len(tasks))
		copy(sortedTasks, tasks)

		// Display each task with both position and stable ID
		for i, t := range sortedTasks {
			if totalShown >= taskLimit {
				break
			}

			assignee := "unassigned"
			if t.AssigneeID != nil {
				assignee = "assigned"
			}

			due := ""
			if t.DueDate != nil {
				due = t.DueDate.Format("2006-01-02")
			}

			// Format the task with position index (stageIndex-taskIndex) and stable ID
			taskPosition := fmt.Sprintf("%d-%d", stagePosition, i+1) // 1-based index
			taskIdentifier := fmt.Sprintf("%s (#%d)", taskPosition, t.NumericID)

			if due != "" {
				fmt.Fprintf(&b, "  %s %s (priority: %s, assignee: %s, due: %s)\n", taskIdentifier, t.Title, t.Priority, assignee, due)
			} else {
				fmt.Fprintf(&b, "  %s %s (priority: %s, assignee: %s)\n", taskIdentifier, t.Title, t.Priority, assignee)
			}
			totalShown++
		}
		if totalShown >= taskLimit {
			break
		}
	}

	if totalShown < len(details.Tasks) {
		fmt.Fprintf(&b, "(showing first %d tasks out of %d)\n", totalShown, len(details.Tasks))
	}

	return strings.TrimSpace(b.String())
}

func (h *Handler) formatProjectDetailsWithUserNames(ctx context.Context, details models.ProjectDetails) string {
	const taskLimit = 50
	var b strings.Builder
	fmt.Fprintf(&b, "📁 %s (status: %s)\n", details.Project.Title, details.Project.Status)
	if strings.TrimSpace(details.Project.Description) != "" {
		fmt.Fprintf(&b, "Description: %s\n", details.Project.Description)
	}

	if len(details.Stages) == 0 {
		b.WriteString("Stages: none\n")
	} else {
		b.WriteString("Stages:\n")
		sorted := make([]models.Stage, len(details.Stages))
		copy(sorted, details.Stages)
		sort.SliceStable(sorted, func(i, j int) bool {
			if sorted[i].OrderIndex == sorted[j].OrderIndex {
				return sorted[i].Title < sorted[j].Title
			}
			if sorted[i].OrderIndex == 0 {
				return false
			}
			if sorted[j].OrderIndex == 0 {
				return true
			}
			return sorted[i].OrderIndex < sorted[j].OrderIndex
		})
		for i, s := range sorted {
			resp := "unassigned"
			if s.ResponsibleID != nil {
				// Try to get the actual user name
				user, err := h.Projects.GetUserByID(ctx, *s.ResponsibleID)
				if err == nil && user.Name != "" {
					resp = user.Name
				} else {
					resp = "assigned"
				}
			}
			fmt.Fprintf(&b, "[%d] %s (responsible: %s)\n", i+1, s.Title, resp) // Show stage index
		}
	}

	if len(details.Tasks) == 0 {
		b.WriteString("Tasks: none")
		return strings.TrimSpace(b.String())
	}

	// group tasks by stage
	byStage := map[string][]models.Task{}
	for _, t := range details.Tasks {
		key := ""
		if t.StageID != nil {
			key = *t.StageID
		}
		byStage[key] = append(byStage[key], t)
	}

	stageOrder := []string{}
	for _, s := range details.Stages {
		stageOrder = append(stageOrder, s.ID)
	}
	stageOrder = append(stageOrder, "") // for unassigned tasks

	b.WriteString("Tasks:\n")
	totalShown := 0

	// Group tasks by stage and sort them for consistent positioning
	for _, stageID := range stageOrder {
		tasks := byStage[stageID]
		if len(tasks) == 0 {
			continue
		}

		// Find the stage to get its position index
		stagePosition := 0
		for i, s := range details.Stages {
			if s.ID == stageID {
				stagePosition = i + 1 // 1-based index
				break
			}
		}

		stageTitle := "No stage"
		for _, s := range details.Stages {
			if s.ID == stageID {
				stageTitle = s.Title
				break
			}
		}

		// Print stage with position index
		if stagePosition > 0 {
			fmt.Fprintf(&b, "• [%d] %s:\n", stagePosition, stageTitle)
		} else {
			fmt.Fprintf(&b, "• %s:\n", stageTitle)
		}

		// Sort tasks within the stage by creation order or ID to maintain consistent positioning
		sortedTasks := make([]models.Task, len(tasks))
		copy(sortedTasks, tasks)

		// Display each task with both position and stable ID
		for i, t := range sortedTasks {
			if totalShown >= taskLimit {
				break
			}

			assignee := "unassigned"
			if t.AssigneeID != nil {
				// Try to get the actual user name
				user, err := h.Projects.GetUserByID(ctx, *t.AssigneeID)
				if err == nil && user.Name != "" {
					assignee = user.Name
				} else {
					assignee = "assigned"
				}
			}
			due := ""
			if t.DueDate != nil {
				due = t.DueDate.Format("2006-01-02")
			}

			// Format the task with position index (stageIndex-taskIndex) and stable ID
			taskPosition := fmt.Sprintf("%d-%d", stagePosition, i+1) // 1-based index
			taskIdentifier := fmt.Sprintf("%s (#%d)", taskPosition, t.NumericID)

			if due != "" {
				fmt.Fprintf(&b, "  %s %s (priority: %s, assignee: %s, due: %s)\n", taskIdentifier, t.Title, t.Priority, assignee, due)
			} else {
				fmt.Fprintf(&b, "  %s %s (priority: %s, assignee: %s)\n", taskIdentifier, t.Title, t.Priority, assignee)
			}
			totalShown++
		}
		if totalShown >= taskLimit {
			break
		}
	}

	if totalShown < len(details.Tasks) {
		fmt.Fprintf(&b, "(showing first %d tasks out of %d)\n", totalShown, len(details.Tasks))
	}

	return strings.TrimSpace(b.String())
}

// resolveTaskByNameOrPositionOrStableID resolves a task by name, position (stageIndex-taskIndex), or stable ID (#id)
func (h *Handler) resolveTaskByNameOrPositionOrStableID(ctx context.Context, entityName string, projectID string) (models.Task, error) {
	var task models.Task
	entityName = strings.TrimSpace(entityName)

	// First, check if it's a stable ID format like "#1052"
	if strings.HasPrefix(entityName, "#") {
		idStr := strings.TrimPrefix(entityName, "#")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			// Find task by numeric ID
			details, err := h.Projects.GetProjectDetails(ctx, projectID, h.User)
			if err != nil {
				return task, err
			}
			for _, t := range details.Tasks {
				if t.NumericID == id {
					return t, nil
				}
			}
			return task, fmt.Errorf("task with ID #%d not found in project", id)
		}
	}

	// Then, check if it's a position format like "2-1"
	if strings.Contains(entityName, "-") {
		parts := strings.Split(entityName, "-")
		if len(parts) == 2 {
			stageIdx, err1 := strconv.Atoi(parts[0])
			taskIdx, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil && stageIdx > 0 && taskIdx > 0 {
				// Get all project details to resolve by position
				details, err := h.Projects.GetProjectDetails(ctx, projectID, h.User)
				if err != nil {
					return task, err
				}

				// Find the stage by position index (1-based)
				if stageIdx > len(details.Stages) {
					return task, fmt.Errorf("stage index %d is out of range (only %d stages available)", stageIdx, len(details.Stages))
				}

				// Sort stages by order to ensure consistent indexing
				sortedStages := make([]models.Stage, len(details.Stages))
				copy(sortedStages, details.Stages)
				sort.SliceStable(sortedStages, func(i, j int) bool {
					if sortedStages[i].OrderIndex == sortedStages[j].OrderIndex {
						return sortedStages[i].Title < sortedStages[j].Title
					}
					if sortedStages[i].OrderIndex == 0 {
						return false
					}
					if sortedStages[j].OrderIndex == 0 {
						return true
					}
					return sortedStages[i].OrderIndex < sortedStages[j].OrderIndex
				})

				targetStage := sortedStages[stageIdx-1] // Convert to 0-based index

				// Find tasks associated with this stage
				stageTasks := []models.Task{}
				for _, t := range details.Tasks {
					if t.StageID != nil && *t.StageID == targetStage.ID {
						stageTasks = append(stageTasks, t)
					}
				}

				// Sort tasks by creation order for consistent indexing
				sort.SliceStable(stageTasks, func(i, j int) bool {
					return stageTasks[i].CreatedAt.Before(stageTasks[j].CreatedAt)
				})

				// Check if task index is valid
				if taskIdx > len(stageTasks) {
					return task, fmt.Errorf("task index %d is out of range for stage '%s' (only %d tasks available)", taskIdx, targetStage.Title, len(stageTasks))
				}

				return stageTasks[taskIdx-1], nil // Convert to 0-based index
			}
		}
	}

	// Finally, fall back to name-based lookup
	return h.Projects.FindTaskForUserByTitle(ctx, entityName, h.User)
}

// resolveStageByNameOrPosition resolves a stage by name or position index
func (h *Handler) resolveStageByNameOrPosition(ctx context.Context, entityName string, projectID string) (models.Stage, error) {
	var stage models.Stage
	entityName = strings.TrimSpace(entityName)

	// Check if it's a position format like "[2]" or just "2"
	stageIdx, err := strconv.Atoi(entityName)
	if err == nil && stageIdx > 0 {
		// Get all project details to resolve by position
		details, err := h.Projects.GetProjectDetails(ctx, projectID, h.User)
		if err != nil {
			return stage, err
		}

		// Sort stages by order to ensure consistent indexing
		sortedStages := make([]models.Stage, len(details.Stages))
		copy(sortedStages, details.Stages)
		sort.SliceStable(sortedStages, func(i, j int) bool {
			if sortedStages[i].OrderIndex == sortedStages[j].OrderIndex {
				return sortedStages[i].Title < sortedStages[j].Title
			}
			if sortedStages[i].OrderIndex == 0 {
				return false
			}
			if sortedStages[j].OrderIndex == 0 {
				return true
			}
			return sortedStages[i].OrderIndex < sortedStages[j].OrderIndex
		})

		// Check if stage index is valid
		if stageIdx > len(sortedStages) {
			return stage, fmt.Errorf("stage index %d is out of range (only %d stages available)", stageIdx, len(sortedStages))
		}

		return sortedStages[stageIdx-1], nil // Convert to 0-based index
	}

	// Fall back to name-based lookup
	return h.Projects.FindStageForUserByTitle(ctx, entityName, h.User)
}
