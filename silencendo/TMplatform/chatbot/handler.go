package chatbot

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

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

func NewHandler(projects *services.ProjectService, user models.User) *Handler {
	return &Handler{Projects: projects, User: user}
}

// Handle attempts to satisfy chatbot intents. handled indicates whether the message was consumed by this layer.
func (h *Handler) Handle(ctx context.Context, message string) (handled bool, reply string, err error) {
	match := DetectIntent(message)
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
		if strings.TrimSpace(match.EntityType) == "" || strings.TrimSpace(match.EntityName) == "" || strings.TrimSpace(match.AssigneeName) == "" {
			return handled, "To assign, say 'Assign Alex to task Permits' or 'Assign Dana to stage Design'.", nil
		}
		if !hasResolved {
			return handled, noActiveProjectMessage(), nil
		}

		assignee, err := h.Projects.EnsureUser(ctx, models.User{Name: strings.TrimSpace(match.AssigneeName)})
		if err != nil {
			return handled, humanizeError(err), err
		}
		switch strings.ToLower(match.EntityType) {
		case "stage":
			stage, err := h.Projects.FindStageForUserByTitle(ctx, match.EntityName, h.User)
			if err != nil {
				return handled, humanizeError(err), err
			}
			if stage.ProjectID != resolvedID {
				return handled, "That stage belongs to another project. Switch active project first.", nil
			}
			if err := h.Projects.AssignResponsible(ctx, "stage", stage.ID, assignee.ID, h.User); err != nil {
				return handled, humanizeError(err), err
			}
			return handled, fmt.Sprintf("✅ Stage '%s' assigned to %s", stage.Title, assignee.Name), nil
		case "task":
			task, err := h.Projects.FindTaskForUserByTitle(ctx, match.EntityName, h.User)
			if err != nil {
				return handled, humanizeError(err), err
			}
			if task.ProjectID != resolvedID {
				return handled, "That task belongs to another project. Switch active project first.", nil
			}
			if err := h.Projects.AssignResponsible(ctx, "task", task.ID, assignee.ID, h.User); err != nil {
				return handled, humanizeError(err), err
			}
			return handled, fmt.Sprintf("✅ Task '%s' assigned to %s", task.Title, assignee.Name), nil
		default:
			return handled, "Please specify whether you want to assign to a stage or a task.", nil
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
		return handled, formatProjectDetails(details), nil
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
		for _, s := range sorted {
			resp := "unassigned"
			if s.ResponsibleID != nil {
				resp = "assigned"
			}
			fmt.Fprintf(&b, "- %s (responsible: %s)\n", s.Title, resp)
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
	for _, stageID := range stageOrder {
		tasks := byStage[stageID]
		if len(tasks) == 0 {
			continue
		}
		stageTitle := "No stage"
		for _, s := range details.Stages {
			if s.ID == stageID {
				stageTitle = s.Title
				break
			}
		}
		fmt.Fprintf(&b, "• %s:\n", stageTitle)
		for _, t := range tasks {
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
			if due != "" {
				fmt.Fprintf(&b, "  - [%s] %s (priority: %s, assignee: %s, due: %s)\n", t.Status, t.Title, t.Priority, assignee, due)
			} else {
				fmt.Fprintf(&b, "  - [%s] %s (priority: %s, assignee: %s)\n", t.Status, t.Title, t.Priority, assignee)
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
