package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"silencendo/chatbot"
	"silencendo/models"
	"silencendo/services"
)

// Dispatcher coordinates routing between project handlers and document/knowledge handlers.
type Dispatcher struct {
	router          *IntentRouter
	projectHandler  *chatbot.Handler
	documentHandler *CLIInterface
	state           *SessionState
}

func NewDispatcher(router *IntentRouter, projectHandler *chatbot.Handler, documentHandler *CLIInterface, state *SessionState) *Dispatcher {
	return &Dispatcher{router: router, projectHandler: projectHandler, documentHandler: documentHandler, state: state}
}

// Dispatch routes a single input using the intent router and appropriate handlers.
// It returns true when the caller should exit the main loop.
func (d *Dispatcher) Dispatch(ctx context.Context, input string) (bool, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return false, nil
	}

	if d.state != nil && d.state.Pending != nil {
		return d.handlePending(ctx, trimmed)
	}

	if d.state != nil && d.state.PendingDeletion != nil {
		handled, err := d.handlePendingDeletion(ctx, trimmed)
		if handled {
			return false, err
		}
	}

	lower := strings.ToLower(trimmed)
	if lower == "/exit" || lower == "/quit" {
		fmt.Println("Goodbye!")
		return true, nil
	}

	if strings.HasPrefix(trimmed, "/") {
		return d.handleSlashCommand(ctx, trimmed)
	}

	if isProjectScopedInput(lower) {
		match := chatbot.DetectIntent(trimmed)
		if match.Intent == chatbot.IntentCreateProject || match.Intent == chatbot.IntentListProjects {
			return d.dispatchWithForced(ctx, trimmed, IntentProject)
		}
		if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) != "" {
			return d.dispatchWithForced(ctx, trimmed, IntentProject)
		}
		fmt.Println(noActiveProjectMessage())
		return false, nil
	}

	resolved := d.resolvedActiveProject()
	route := d.router.Route(trimmed, resolved)
	switch route.Type {
	case IntentAmbiguous:
		if d.state != nil {
			d.state.Pending = &PendingConfirmation{Kind: ConfirmationProjectVsDocument, OriginalInput: trimmed, CreatedAt: time.Now()}
		}
		fmt.Println("Do you want to:\n1) Create a new PROJECT in the system\n2) Edit the CURRENT document?")
		return false, nil
	case IntentProject:
		return d.dispatchWithForced(ctx, trimmed, IntentProject)
	case IntentDocument:
		return d.dispatchWithForced(ctx, trimmed, IntentDocument)
	case IntentKnowledge:
		return false, d.documentHandler.handleQuestion(trimmed)
	default:
		return false, nil
	}
}

func (d *Dispatcher) handlePending(ctx context.Context, input string) (bool, error) {
	if d.state == nil || d.state.Pending == nil {
		return false, nil
	}
	lower := strings.ToLower(strings.TrimSpace(input))

	projectChoices := map[string]bool{"1": true, "project": true, "create project": true, "new project": true, "yes": true}
	documentChoices := map[string]bool{"2": true, "document": true, "edit document": true, "file": true}
	cancelChoices := map[string]bool{"cancel": true, "exit": true, "quit": true, "no": true}

	if projectChoices[lower] {
		original := d.state.Pending.OriginalInput
		d.state.Pending = nil
		return d.dispatchWithForced(ctx, original, IntentProject)
	}

	if documentChoices[lower] {
		original := d.state.Pending.OriginalInput
		d.state.Pending = nil
		return d.dispatchWithForced(ctx, original, IntentDocument)
	}

	if cancelChoices[lower] {
		d.state.Pending = nil
		fmt.Println("Cancelled.")
		return false, nil
	}

	fmt.Println("Please reply with 1 (project), 2 (document), or 'cancel'.")
	return false, nil
}

func (d *Dispatcher) handlePendingDeletion(ctx context.Context, input string) (bool, error) {
	if d.state == nil || d.state.PendingDeletion == nil {
		return false, nil
	}

	title := strings.TrimSpace(d.state.PendingDeletion.ProjectTitle)
	lower := strings.ToLower(strings.TrimSpace(input))
	if lower == "cancel" {
		d.state.PendingDeletion = nil
		fmt.Println("Deletion cancelled.")
		return true, nil
	}

	expected := fmt.Sprintf("delete %s", strings.ToLower(title))
	if lower == expected {
		projID := d.state.PendingDeletion.ProjectID
		d.state.PendingDeletion = nil
		if projID == "" {
			return true, nil
		}

		if d.projectHandler == nil || d.projectHandler.Projects == nil {
			return true, fmt.Errorf("project service unavailable")
		}

		deleted, err := d.projectHandler.Projects.SoftDeleteProject(ctx, projID, d.projectHandler.User)
		if err != nil {
			return true, err
		}
		if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) == projID {
			d.state.ActiveProjectID = ""
			d.state.ActiveProjectTitle = ""
		}
		fmt.Printf("🗑️ Project '%s' deleted.\n", deleted.Title)
		return true, nil
	}

	fmt.Printf("To confirm deletion, reply with 'delete %s' or 'cancel'.\n", title)
	return true, nil
}

func (d *Dispatcher) dispatchWithForced(ctx context.Context, input string, forced IntentType) (bool, error) {
	switch forced {
	case IntentProject:
		if d.projectHandler != nil {
			if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) != "" {
				ctx = chatbot.WithResolvedProject(ctx, strings.TrimSpace(d.state.ActiveProjectID), strings.TrimSpace(d.state.ActiveProjectTitle))
			}
			handled, reply, err := d.projectHandler.Handle(ctx, input)
			if handled {
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else if strings.TrimSpace(reply) != "" {
					fmt.Println(reply)
				}
				return false, err
			}
			if err != nil {
				return false, err
			}
		}
		// If not handled, avoid falling back to LLM for project intents; give direct guidance.
		if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) == "" {
			fmt.Println(noActiveProjectMessage())
			return false, nil
		}
		fmt.Println("I couldn't understand that project request. Try /project list, /project use <number>, or /project current.")
		return false, nil
	case IntentDocument:
		return false, d.documentHandler.handleEditRequest(input)
	default:
		return false, nil
	}
}

func (d *Dispatcher) handleSlashCommand(ctx context.Context, input string) (bool, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false, nil
	}

	command := strings.TrimPrefix(strings.ToLower(parts[0]), "/")
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}

	if command == "project" || command == "projects" || command == "проект" || command == "проекты" {
		return false, d.handleProjectSlashCommand(ctx, args)
	}

	if err := d.documentHandler.handleCommand(input); err != nil {
		return false, err
	}
	return false, nil
}

func (d *Dispatcher) resolvedActiveProject() *ResolvedProject {
	if d.state == nil {
		return nil
	}
	if strings.TrimSpace(d.state.ActiveProjectID) == "" {
		return nil
	}
	return &ResolvedProject{ID: strings.TrimSpace(d.state.ActiveProjectID), Title: strings.TrimSpace(d.state.ActiveProjectTitle), Confidence: 1.0}
}

func (d *Dispatcher) handleProjectSlashCommand(ctx context.Context, args []string) error {
	svc := d.projectHandler.Projects
	if svc == nil {
		return fmt.Errorf("project service unavailable")
	}

	if len(args) == 0 {
		return d.printProjectList(ctx, false)
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "list", "список", "list_projects", "listprojects", "мои", "мои_проекты", "проекты", "projects":
		verbose := containsFlag(args[1:], "--verbose")
		return d.printProjectList(ctx, verbose)
	case "use", "set", "select", "выбрать", "использовать", "выбор":
		if len(args) < 2 {
			fmt.Println("Usage: /project use <name|id|number>")
			return nil
		}
		target := strings.Join(args[1:], " ")
		return d.selectActiveProject(ctx, target)
	case "current", "active", "текущий", "активный":
		if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) != "" {
			fmt.Println(d.describeCurrentProject())
			return nil
		}
		fmt.Println(noActiveProjectMessage())
		return nil
	case "clear", "сброс", "очистить":
		if d.state != nil {
			d.state.ActiveProjectID = ""
			d.state.ActiveProjectTitle = ""
		}
		fmt.Println("Active project cleared.")
		return nil
	case "set-status", "status", "setstatus", "статус":
		if len(args) < 3 {
			fmt.Println("Usage: /project set-status <active|inactive> <name|id|number>")
			return nil
		}
		status := strings.ToLower(strings.TrimSpace(args[1]))
		target := strings.Join(args[2:], " ")
		return d.updateProjectStatus(ctx, status, target)
	case "delete", "remove", "удали", "удалить":
		if len(args) < 2 {
			fmt.Println("Usage: /project delete <name|id|number>")
			return nil
		}
		target := strings.Join(args[1:], " ")
		return d.requestProjectDeletion(ctx, target)
	default:
		fmt.Println("Unknown /project subcommand. Try: list, use, current, clear, set-status, delete.")
		return nil
	}
}

func (d *Dispatcher) printProjectList(ctx context.Context, verbose bool) error {
	projects, err := d.projectHandler.Projects.ListProjects(ctx, d.projectHandler.User, services.ListFilters{})
	if err != nil {
		return err
	}

	if d.state != nil {
		d.state.LastProjectList = []ProjectSummary{}
	}

	if len(projects) == 0 {
		fmt.Println("You have no projects yet. Create one with 'create project <name>'.")
		return nil
	}

	activeID := ""
	if d.state != nil {
		activeID = strings.TrimSpace(d.state.ActiveProjectID)
	}

	var lines []string
	for i, p := range projects {
		marker := "  "
		if p.ID == activeID {
			marker = "* "
		}
		line := fmt.Sprintf("%s[%d] %s (status: %s)", marker, i+1, p.Title, p.Status)
		if verbose {
			line = fmt.Sprintf("%s id: %s", line, p.ID)
		}
		lines = append(lines, line)
		if d.state != nil {
			d.state.LastProjectList = append(d.state.LastProjectList, ProjectSummary{ID: p.ID, Title: p.Title, Status: p.Status, NormalizedTitle: strings.ToLower(strings.TrimSpace(p.Title))})
		}
	}

	fmt.Println("📁 Your projects:")
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}

func (d *Dispatcher) selectActiveProject(ctx context.Context, target string) error {
	project, err := d.resolveProjectSelection(ctx, target)
	if err != nil {
		return err
	}
	if project.ID == "" {
		return nil
	}

	previousID := ""
	if d.state != nil {
		previousID = strings.TrimSpace(d.state.ActiveProjectID)
	}

	activated, deactivated, err := d.projectHandler.Projects.SwitchActiveProject(ctx, project.ID, previousID, d.projectHandler.User)
	if err != nil {
		return err
	}

	if d.state != nil {
		d.state.ActiveProjectID = activated.ID
		d.state.ActiveProjectTitle = activated.Title
	}

	msg := fmt.Sprintf("✅ Current project set to: %s", activated.Title)
	if deactivated != nil {
		msg = fmt.Sprintf("%s (previous project '%s' set to inactive)", msg, deactivated.Title)
	}
	fmt.Println(msg)
	return nil
}

func (d *Dispatcher) updateProjectStatus(ctx context.Context, status, target string) error {
	if status != "active" && status != "inactive" && status != "активный" && status != "неактивный" {
		fmt.Println("Status must be 'active' or 'inactive'.")
		return nil
	}

	if status == "активный" {
		status = "active"
	}
	if status == "неактивный" {
		status = "inactive"
	}

	project, err := d.resolveProjectSelection(ctx, target)
	if err != nil {
		return err
	}
	if project.ID == "" {
		return nil
	}

	if status == "active" {
		previousID := ""
		if d.state != nil {
			previousID = strings.TrimSpace(d.state.ActiveProjectID)
		}
		updated, deactivated, err := d.projectHandler.Projects.SwitchActiveProject(ctx, project.ID, previousID, d.projectHandler.User)
		if err != nil {
			return err
		}
		if d.state != nil {
			d.state.ActiveProjectID = updated.ID
			d.state.ActiveProjectTitle = updated.Title
		}
		msg := fmt.Sprintf("Project %s status set to active and selected as current.", updated.Title)
		if deactivated != nil {
			msg = fmt.Sprintf("%s Previous project '%s' set to inactive.", msg, deactivated.Title)
		}
		fmt.Println(msg)
		return nil
	}

	updated, err := d.projectHandler.Projects.SetProjectStatus(ctx, project.ID, status, d.projectHandler.User)
	if err != nil {
		return err
	}

	if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) == updated.ID {
		d.state.ActiveProjectID = ""
		d.state.ActiveProjectTitle = ""
		fmt.Printf("Project %s marked inactive and current selection cleared.\n", updated.Title)
		return nil
	}

	fmt.Printf("Project %s status set to %s. Current session project remains %s.\n", updated.Title, status, d.currentProjectLabel())
	return nil
}

func (d *Dispatcher) requestProjectDeletion(ctx context.Context, target string) error {
	project, err := d.resolveProjectSelection(ctx, target)
	if err != nil {
		return err
	}
	if project.ID == "" {
		return nil
	}

	if d.state != nil {
		d.state.PendingDeletion = &PendingDeletion{ProjectID: project.ID, ProjectTitle: project.Title}
	}

	fmt.Printf("⚠️ Delete project '%s'? Reply 'delete %s' to confirm or 'cancel'.\n", project.Title, project.Title)
	return nil
}

func (d *Dispatcher) resolveProjectSelection(ctx context.Context, target string) (models.Project, error) {
	if d.projectHandler == nil || d.projectHandler.Projects == nil {
		return models.Project{}, fmt.Errorf("project service unavailable")
	}

	projects, err := d.projectHandler.Projects.ListProjects(ctx, d.projectHandler.User, services.ListFilters{})
	if err != nil {
		return models.Project{}, err
	}

	normalizedTarget := strings.ToLower(strings.TrimSpace(target))
	if normalizedTarget == "" {
		fmt.Println("Please provide a project name, number, or id prefix.")
		return models.Project{}, nil
	}

	// 1) numeric index against cached list
	if idx, errConv := strconv.Atoi(normalizedTarget); errConv == nil {
		if d.state == nil || len(d.state.LastProjectList) == 0 {
			fmt.Println("No cached project list. Run /project list first.")
			return models.Project{}, nil
		}
		if idx <= 0 || idx > len(d.state.LastProjectList) {
			fmt.Println("Invalid project number. Use /project list and pick a listed number.")
			return models.Project{}, nil
		}
		selected := d.state.LastProjectList[idx-1]
		for _, p := range projects {
			if p.ID == selected.ID {
				return p, nil
			}
		}
	}

	// 2) exact normalized title
	var titleMatches []models.Project
	for _, p := range projects {
		if strings.ToLower(strings.TrimSpace(p.Title)) == normalizedTarget {
			titleMatches = append(titleMatches, p)
		}
	}
	if len(titleMatches) == 1 {
		return titleMatches[0], nil
	}
	if len(titleMatches) > 1 {
		fmt.Println("Multiple projects match that title. Use /project list and pick a number.")
		return models.Project{}, nil
	}

	// 3) UUID prefix
	var prefixMatches []models.Project
	for _, p := range projects {
		if strings.HasPrefix(strings.ToLower(p.ID), normalizedTarget) {
			prefixMatches = append(prefixMatches, p)
		}
	}
	if len(prefixMatches) == 1 {
		return prefixMatches[0], nil
	}
	if len(prefixMatches) > 1 {
		fmt.Println("Multiple projects match that id prefix. Use /project list --verbose to choose a number.")
		return models.Project{}, nil
	}

	fmt.Println("Project not found. Use /project list to see available projects.")
	return models.Project{}, nil
}

func noActiveProjectMessage() string {
	return "No active project selected. Use /project list and /project use <number>."
}

func (d *Dispatcher) describeCurrentProject() string {
	label := d.currentProjectLabel()
	if label == "none" {
		return "Active project: none"
	}
	return fmt.Sprintf("Active project: %s", label)
}

func (d *Dispatcher) currentProjectLabel() string {
	if d.state == nil || strings.TrimSpace(d.state.ActiveProjectID) == "" {
		return "none"
	}
	activeID := strings.TrimSpace(d.state.ActiveProjectID)
	activeTitle := strings.TrimSpace(d.state.ActiveProjectTitle)
	if len(d.state.LastProjectList) > 0 {
		for idx, p := range d.state.LastProjectList {
			if p.ID == activeID {
				return fmt.Sprintf("[%d] %s", idx+1, activeTitle)
			}
		}
	}
	return activeTitle
}

func containsFlag(args []string, flag string) bool {
	for _, a := range args {
		if strings.TrimSpace(a) == flag {
			return true
		}
	}
	return false
}

func isProjectScopedInput(lower string) bool {
	keywords := []string{"project", "plan", "task", "stage", "member", "assign", "status", "progress", "delete", "проект", "план", "задача", "этап", "участник", "назнач", "статус", "удал", "прогресс"}
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
