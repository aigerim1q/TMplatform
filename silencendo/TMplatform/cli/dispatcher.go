package cli

import (
	"context"
	"fmt"
	"regexp"
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

	if d.state != nil && d.state.PendingAssign != nil {
		return d.handlePendingAssignment(ctx, trimmed)
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

	if (d.state == nil || strings.TrimSpace(d.state.ActiveProjectID) == "") && isProjectScopedInput(lower) {
		match := chatbot.DetectIntent(trimmed)
		if match.Intent == chatbot.IntentCreateProject || match.Intent == chatbot.IntentListProjects {
			return d.dispatchWithForced(ctx, trimmed, IntentProject, nil)
		}
		fmt.Println(noActiveProjectMessage())
		return false, nil
	}

	resolved := d.resolvedActiveProject()
	route := d.router.Route(ctx, trimmed, resolved)
	switch route.Type {
	case IntentAmbiguous:
		if d.state != nil {
			d.state.Pending = &PendingConfirmation{Kind: ConfirmationProjectVsDocument, OriginalInput: trimmed, CreatedAt: time.Now()}
		}
		fmt.Println("Do you want to:\n1) Create a new PROJECT in the system\n2) Edit the CURRENT document?")
		return false, nil
	case IntentProject:
		return d.dispatchWithForced(ctx, trimmed, IntentProject, route.ProjectMatch)
	case IntentDocument:
		return d.dispatchWithForced(ctx, trimmed, IntentDocument, nil)
	case IntentKnowledge:
		if route.LLMFailed {
			fmt.Println("DeepSeek couldn't understand that request. Please rephrase.")
			return false, nil
		}
		return false, d.documentHandler.handleQuestion(trimmed)
	default:
		return false, nil
	}
}

func (d *Dispatcher) handlePending(ctx context.Context, input string) (bool, error) {
	if d.state == nil || d.state.Pending == nil {
		// Check for pending assignment if no other pending state exists
		if d.state != nil && d.state.PendingAssign != nil {
			return d.handlePendingAssignment(ctx, input)
		}
		return false, nil
	}
	lower := strings.ToLower(strings.TrimSpace(input))

	projectChoices := map[string]bool{"1": true, "project": true, "create project": true, "new project": true, "yes": true}
	documentChoices := map[string]bool{"2": true, "document": true, "edit document": true, "file": true}
	cancelChoices := map[string]bool{"cancel": true, "exit": true, "quit": true, "no": true}

	if projectChoices[lower] {
		original := d.state.Pending.OriginalInput
		d.state.Pending = nil
		return d.dispatchWithForced(ctx, original, IntentProject, nil)
	}

	if documentChoices[lower] {
		original := d.state.Pending.OriginalInput
		d.state.Pending = nil
		return d.dispatchWithForced(ctx, original, IntentDocument, nil)
	}

	if cancelChoices[lower] {
		d.state.Pending = nil
		fmt.Println("Cancelled.")
		return false, nil
	}

	fmt.Println("Please reply with 1 (project), 2 (document), or 'cancel'.")
	return false, nil
}

// handlePendingAssignment handles responses to assignment clarification questions
func (d *Dispatcher) handlePendingAssignment(ctx context.Context, input string) (bool, error) {
	if d.state == nil || d.state.PendingAssign == nil {
		return false, nil
	}

	pending := d.state.PendingAssign

	// Parse the input for task references and assignees
	taskStableID, taskPosStage, taskPosIndex := parseTaskReference(input)

	// Parse assignees if needed
	var assignees []string
	if pending.WaitingFor == "assignees" {
		assignees = parseAssigneesFromInput(input)
	} else {
		assignees = pending.Assignees
	}

	// Update the pending assignment with parsed data
	if taskStableID != 0 {
		pending.TaskStableID = taskStableID
	}
	if taskPosStage != 0 && taskPosIndex != 0 {
		pending.TaskPosStage = taskPosStage
		pending.TaskPosIndex = taskPosIndex
	}
	if len(assignees) > 0 {
		pending.Assignees = assignees
	}

	// Check if the assignment is now complete
	if (len(pending.Assignees) > 0 && (pending.TaskStableID != 0 || (pending.TaskPosStage != 0 && pending.TaskPosIndex != 0))) ||
		(pending.WaitingFor == "assignees" && len(pending.Assignees) > 0) ||
		(pending.WaitingFor == "task_ref" && (pending.TaskStableID != 0 || (pending.TaskPosStage != 0 && pending.TaskPosIndex != 0))) {

		// Execute the assignment
		err := d.executePendingAssignment(ctx, pending)
		if err != nil {
			fmt.Printf("Error executing assignment: %v\n", err)
			d.state.PendingAssign = nil
			return false, err
		}

		// Clear the pending assignment
		d.state.PendingAssign = nil
		return false, nil
	}

	// If still not complete, ask for the missing piece
	if pending.WaitingFor == "assignees" {
		if pending.TaskStableID != 0 {
			fmt.Printf("Who should I assign to task #%d?\n", pending.TaskStableID)
		} else if pending.TaskPosStage != 0 && pending.TaskPosIndex != 0 {
			fmt.Printf("Who should I assign to task %d-%d?\n", pending.TaskPosStage, pending.TaskPosIndex)
		} else {
			fmt.Println("Who should I assign to which task?")
		}
	} else { // waiting for task reference
		fmt.Println("Which task should I assign them to? Use task 4-2 or #8.")
	}

	return false, nil
}

// parseTaskReference extracts task references from input
func parseTaskReference(input string) (int64, int, int) {
	// Match stable ID: #\s*(\d+)
	stableIDRegex := regexp.MustCompile(`#\s*(\d+)`)
	stableMatches := stableIDRegex.FindStringSubmatch(input)
	if len(stableMatches) >= 2 {
		if id, err := strconv.ParseInt(stableMatches[1], 10, 64); err == nil {
			return id, 0, 0
		}
	}

	// Match positional: (\d+)\s*-\s*(\d+)
	posRegex := regexp.MustCompile(`(\d+)\s*-\s*(\d+)`)
	posMatches := posRegex.FindStringSubmatch(input)
	if len(posMatches) >= 3 {
		if stage, err := strconv.Atoi(posMatches[1]); err == nil {
			if index, err := strconv.Atoi(posMatches[2]); err == nil {
				return 0, stage, index
			}
		}
	}

	return 0, 0, 0
}

// parseAssigneesFromInput parses assignees from input text
func parseAssigneesFromInput(input string) []string {
	text := strings.TrimSpace(input)

	// Remove common words that aren't assignees
	// Check if input starts with common assignment words and remove them
	lowerText := strings.ToLower(text)
	if strings.HasPrefix(lowerText, "assign") || strings.HasPrefix(lowerText, "make") ||
		strings.HasPrefix(lowerText, "let") || strings.HasPrefix(lowerText, "put") ||
		strings.HasPrefix(lowerText, "get") || strings.HasPrefix(lowerText, "have") ||
		strings.HasPrefix(lowerText, "add") {
		// Extract assignees after the verb
		parts := strings.SplitN(lowerText, " ", 2)
		if len(parts) > 1 {
			text = parts[1]
		}
	}

	// Split by "and" and "," to get multiple assignees
	text = strings.ReplaceAll(text, " and ", ",")
	text = strings.ReplaceAll(text, " & ", ",")
	text = strings.ReplaceAll(text, " + ", ",")

	assigneeParts := strings.FieldsFunc(text, func(r rune) bool {
		return r == ','
	})

	var assignees []string
	for _, part := range assigneeParts {
		cleanPart := strings.TrimSpace(part)
		// Filter out empty strings and common non-name words
		if cleanPart != "" && cleanPart != "to" && cleanPart != "the" && cleanPart != "a" && cleanPart != "an" {
			assignees = append(assignees, cleanPart)
		}
	}

	return assignees
}

// executePendingAssignment executes the assignment based on the pending assignment data
func (d *Dispatcher) executePendingAssignment(ctx context.Context, pending *PendingAssign) error {
	if d.projectHandler == nil {
		return fmt.Errorf("project handler not available")
	}

	// Build a fake input string to pass to the handler for processing
	var entityName string
	if pending.TaskStableID != 0 {
		entityName = fmt.Sprintf("#%d", pending.TaskStableID)
	} else if pending.TaskPosStage != 0 && pending.TaskPosIndex != 0 {
		entityName = fmt.Sprintf("%d-%d", pending.TaskPosStage, pending.TaskPosIndex)
	} else {
		return fmt.Errorf("no task reference provided")
	}

	// Combine assignees into a single string
	assigneeStr := strings.Join(pending.Assignees, " and ")

	// Create a temporary context with the resolved project
	ctxWithProject := chatbot.WithResolvedProject(ctx, pending.ProjectID, pending.ProjectTitle)

	// Create an intent match for the assignment
	match := chatbot.IntentMatch{
		Intent:       chatbot.IntentAssignResponsible,
		EntityType:   "task",
		EntityName:   entityName,
		AssigneeName: assigneeStr,
		Confidence:   0.9,
	}

	// Execute the assignment via the project handler
	handled, reply, err := d.projectHandler.HandleWithMatch(ctxWithProject, match)
	if err != nil {
		return err
	}

	if handled && reply != "" {
		fmt.Println(reply)
	}

	return nil
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
			d.state.ClearActiveProject()
		}
		fmt.Printf("🗑️ Project '%s' deleted.\n", deleted.Title)
		return true, nil
	}

	fmt.Printf("To confirm deletion, reply with 'delete %s' or 'cancel'.\n", title)
	return true, nil
}

func (d *Dispatcher) dispatchWithForced(ctx context.Context, input string, forced IntentType, projectMatch *chatbot.IntentMatch) (bool, error) {
	switch forced {
	case IntentProject:
		if d.projectHandler != nil {
			if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) != "" {
				ctx = chatbot.WithResolvedProject(ctx, strings.TrimSpace(d.state.ActiveProjectID), strings.TrimSpace(d.state.ActiveProjectTitle))
			}
			var handled bool
			var reply string
			var err error
			if projectMatch != nil {
				// Check if this is an assignment that needs to be made pending
				if projectMatch.Intent == chatbot.IntentAssignResponsible {
					handled, reply, err = d.handleAssignmentWithPartialInfo(ctx, input, *projectMatch)
				} else {
					handled, reply, err = d.projectHandler.HandleWithMatch(ctx, *projectMatch)
				}
			} else {
				handled, reply, err = d.projectHandler.Handle(ctx, input)

				// Check if the handler returned a partial assignment request
				isPartialAssignment := reply == "Which task? Use task 4-2 or #8." ||
					reply == "Which task should I assign them to? Use task 4-2 or #8." ||
					reply == "Who should I assign to task #?" ||
					strings.Contains(reply, "Who should I assign to task #") ||
					strings.Contains(reply, "Could you clarify who you want to assign to what task or stage") ||
					strings.Contains(reply, "Which task should I assign them to") ||
					strings.Contains(reply, "who you want to assign to what task or stage") ||
					(strings.Contains(strings.ToLower(input), "assign") && strings.Contains(strings.ToLower(input), "#") && strings.Contains(reply, "clarify who you want to assign"))

				if handled && isPartialAssignment {

					// Parse the input to extract assignees
					assignees := parseAssigneesFromInput(input)
					if len(assignees) > 0 {
						// Create a pending assignment for task reference
						if d.state != nil {
							id, title, ok := chatbot.ResolvedProjectFromCtx(ctx)
							if ok {
								d.state.PendingAssign = &PendingAssign{
									Assignees:    assignees,
									WaitingFor:   "task_ref",
									ProjectID:    id,
									ProjectTitle: title,
								}
							}
						}
						fmt.Println(reply)
						return false, nil
					}
				}
			}
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

			// Fallback: with an active project, show its details when free-form text isn't parsed.
			if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) != "" {
				fallbackHandled, fallbackReply, fallbackErr := d.projectHandler.Handle(ctx, "show project")
				if fallbackHandled {
					if fallbackErr != nil {
						fmt.Printf("Error: %v\n", fallbackErr)
					} else if strings.TrimSpace(fallbackReply) != "" {
						fmt.Println(fallbackReply)
					}
					return false, fallbackErr
				}
				if fallbackErr != nil {
					return false, fallbackErr
				}
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

// handleAssignmentWithPartialInfo checks if the assignment has partial information and creates a pending assignment if needed
func (d *Dispatcher) handleAssignmentWithPartialInfo(ctx context.Context, input string, match chatbot.IntentMatch) (bool, string, error) {
	entityName := strings.TrimSpace(match.EntityName)
	assigneeName := strings.TrimSpace(match.AssigneeName)

	// Check if we have partial information that requires clarification
	if entityName == "" || assigneeName == "" {
		if d.state != nil {
			id, title, ok := chatbot.ResolvedProjectFromCtx(ctx)
			if ok {
				// Determine what information is missing
				if entityName == "" && assigneeName != "" {
					// Missing task reference but have assignees - split assignees if needed
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

					// Create pending assignment waiting for task reference
					d.state.PendingAssign = &PendingAssign{
						Assignees:    assigneeNames,
						WaitingFor:   "task_ref",
						ProjectID:    id,
						ProjectTitle: title,
					}

					return true, "Which task should I assign them to? Use task 4-2 or #8.", nil
				} else if assigneeName == "" && entityName != "" {
					// Missing assignees but have task reference
					var stableID int64
					var posStage, posIndex int

					// Try to parse the entity name for task reference
					stableIDRegex := regexp.MustCompile(`#\s*(\d+)`)
					stableMatches := stableIDRegex.FindStringSubmatch(entityName)
					if len(stableMatches) >= 2 {
						if id, err := strconv.ParseInt(stableMatches[1], 10, 64); err == nil {
							stableID = id
						}
					}

					// Also check for positional format
					posRegex := regexp.MustCompile(`(\d+)\s*-\s*(\d+)`)
					posMatches := posRegex.FindStringSubmatch(entityName)
					if len(posMatches) >= 3 {
						if stage, err := strconv.Atoi(posMatches[1]); err == nil {
							if index, err := strconv.Atoi(posMatches[2]); err == nil {
								posStage = stage
								posIndex = index
							}
						}
					}

					// Create pending assignment waiting for assignees
					d.state.PendingAssign = &PendingAssign{
						TaskStableID: stableID,
						TaskPosStage: posStage,
						TaskPosIndex: posIndex,
						WaitingFor:   "assignees",
						ProjectID:    id,
						ProjectTitle: title,
					}

					// Generate appropriate message based on task reference type
					if stableID != 0 {
						return true, fmt.Sprintf("Who should I assign to task #%d?", stableID), nil
					} else if posStage != 0 && posIndex != 0 {
						return true, fmt.Sprintf("Who should I assign to task %d-%d?", posStage, posIndex), nil
					} else {
						return true, "Who should I assign to which task?", nil
					}
				}
			}
		}
	}

	// If we have complete information, just delegate to the handler
	return d.projectHandler.HandleWithMatch(ctx, match)
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
	case "help", "h", "?":
		d.printProjectHelp()
		return nil
	case "use", "set", "select", "выбрать", "использовать", "выбор":
		if len(args) < 2 {
			fmt.Println("Usage: /project use <name|id|number>")
			return nil
		}
		target := strings.Join(args[1:], " ")
		return d.selectActiveProject(ctx, target)
	case "current", "active", "текущий", "активный":
		if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) != "" {
			// Instead of just showing the project name, show the full project details like show_project intent
			ctxWithProject := chatbot.WithResolvedProject(ctx, strings.TrimSpace(d.state.ActiveProjectID), strings.TrimSpace(d.state.ActiveProjectTitle))
			handled, reply, err := d.projectHandler.Handle(ctxWithProject, "show project")
			if handled {
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else if strings.TrimSpace(reply) != "" {
					fmt.Println(reply)
				}
			} else {
				// Fallback to just showing the project name
				fmt.Println(d.describeCurrentProject())
			}
			return nil
		}
		fmt.Println(noActiveProjectMessage())
		return nil
	case "clear", "сброс", "очистить":
		if d.state != nil {
			d.state.ClearActiveProject()
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
		fmt.Println("Unknown /project subcommand. Try: help, list, use, current, clear, set-status, delete.")
		return nil
	}
}

func (d *Dispatcher) printProjectHelp() {
	fmt.Println("/project help - Show this help")
	fmt.Println("/project list [--verbose] - List your projects")
	fmt.Println("/project use <name|id|number> - Set active project")
	fmt.Println("/project current - Show active project")
	fmt.Println("/project clear - Clear active project")
	fmt.Println("/project set-status <active|inactive> <name|id|number> - Update status")
	fmt.Println("/project delete <name|id|number> - Soft delete a project (confirm required)")
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
		activeID := ""
		activeTitle := ""
		if d.state != nil {
			activeID = strings.TrimSpace(d.state.ActiveProjectID)
			activeTitle = strings.TrimSpace(d.state.ActiveProjectTitle)
		}

		if activeID != "" {
			project = models.Project{ID: activeID, Title: activeTitle}
		} else if d.projectHandler != nil && d.projectHandler.Projects != nil {
			projects, err := d.projectHandler.Projects.ListProjects(ctx, d.projectHandler.User, services.ListFilters{})
			if err == nil && len(projects) == 1 {
				project = projects[0]
			}
		}

		if project.ID == "" {
			return nil
		}
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
		d.state.SetActiveProject(activated.ID, activated.Title)
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
			d.state.SetActiveProject(updated.ID, updated.Title)
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
		d.state.ClearActiveProject()
		fmt.Printf("Project %s marked inactive and current selection cleared.\n", updated.Title)
		return nil
	}

	fmt.Printf("Project %s status set to %s. Current session project remains %s.\n", updated.Title, status, d.currentProjectLabel())
	return nil
}

func (d *Dispatcher) requestProjectDeletion(ctx context.Context, target string) error {
	if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) != "" {
		project := models.Project{ID: strings.TrimSpace(d.state.ActiveProjectID), Title: strings.TrimSpace(d.state.ActiveProjectTitle)}
		if d.state != nil {
			d.state.PendingDeletion = &PendingDeletion{ProjectID: project.ID, ProjectTitle: project.Title}
		}
		fmt.Printf("⚠️ Delete project '%s'? Reply 'delete %s' to confirm or 'cancel'.\n", project.Title, project.Title)
		return nil
	}

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
		if idx <= 0 {
			fmt.Println("Invalid project number. Use /project list and pick a listed number.")
			return models.Project{}, nil
		}
		if idx <= len(projects) {
			return projects[idx-1], nil
		}

		if d.state != nil && strings.TrimSpace(d.state.ActiveProjectID) != "" {
			activeID := strings.TrimSpace(d.state.ActiveProjectID)
			activeTitle := strings.TrimSpace(d.state.ActiveProjectTitle)
			for _, p := range projects {
				if p.ID == activeID {
					return p, nil
				}
			}
			return models.Project{ID: activeID, Title: activeTitle}, nil
		}

		fmt.Println("Invalid project number. Use /project list and pick a listed number.")
		return models.Project{}, nil
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
