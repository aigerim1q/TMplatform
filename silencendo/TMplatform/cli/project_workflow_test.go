package cli

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"silencendo/models"
	"silencendo/services"
)

func TestProjectListShowsNumericIndicesOnly(t *testing.T) {
	ctx := context.Background()
	dbConn := connectTestDB(t, "file:project-list-numeric.db?_foreign_keys=1")
	user := models.User{ID: "user-list", Name: "Lister"}
	dispatcher, _, svc := newTestDispatcher(t, dbConn, user)

	_, _, _ = svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "minecraft speedrun"}, user)
	_, _, _ = svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "making panini"}, user)

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(ctx, "/project list")
	})

	if strings.Contains(output, "id:") {
		t.Fatalf("list output should not show ids: %s", output)
	}
	if !strings.Contains(output, "[1]") || !strings.Contains(output, "[2]") {
		t.Fatalf("expected numeric indices, got: %s", output)
	}
	if !strings.Contains(output, "minecraft speedrun") || !strings.Contains(output, "making panini") {
		t.Fatalf("expected titles present, got: %s", output)
	}
}

func TestProjectUseUpdatesStatuses(t *testing.T) {
	ctx := context.Background()
	dbConn := connectTestDB(t, "file:project-use-switch.db?_foreign_keys=1")
	user := models.User{ID: "user-active", Name: "Active User"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, user)

	first, _, _ := svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "minecraft speedrun"}, user)
	second, _, _ := svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "making panini"}, user)

	_, _ = dispatcher.Dispatch(ctx, "/project list")
	if len(state.LastProjectList) != 2 {
		t.Fatalf("expected 2 cached projects, got %d", len(state.LastProjectList))
	}
	firstListed := state.LastProjectList[0]
	secondListed := state.LastProjectList[1]

	if _, err := dispatcher.Dispatch(ctx, "/project use 1"); err != nil {
		t.Fatalf("use first project: %v", err)
	}
	if state.ActiveProjectID != firstListed.ID {
		t.Fatalf("active project mismatch, want %s got %s", firstListed.ID, state.ActiveProjectID)
	}

	if _, err := dispatcher.Dispatch(ctx, "/project use 2"); err != nil {
		t.Fatalf("use second project: %v", err)
	}
	if state.ActiveProjectID != secondListed.ID {
		t.Fatalf("active project should switch to second selection, got %s", state.ActiveProjectID)
	}

	projects, err := svc.ListProjects(ctx, user, services.ListFilters{})
	if err != nil {
		t.Fatalf("list projects after switches: %v", err)
	}
	status := map[string]string{}
	for _, p := range projects {
		status[p.ID] = p.Status
	}
	if status[firstListed.ID] != "inactive" {
		t.Fatalf("first listed project should become inactive, got %s", status[firstListed.ID])
	}
	if status[secondListed.ID] != "active" {
		t.Fatalf("second listed project should be active, got %s", status[secondListed.ID])
	}
}

func TestPlanCreationUsesActiveProject(t *testing.T) {
	ctx := context.Background()
	dbConn := connectTestDB(t, "file:plan-active.db?_foreign_keys=1")
	user := models.User{ID: "user-route", Name: "Route User"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, user)

	proj, _, _ := svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "minecraft speedrun"}, user)
	_, _ = dispatcher.Dispatch(ctx, "/project list")
	_, _ = dispatcher.Dispatch(ctx, "/project use 1")
	if state.ActiveProjectID == "" {
		t.Fatalf("expected active project to be set")
	}

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(ctx, "create plan for Yusuf and Omar to beat minecraft")
	})
	if strings.Contains(output, "Do you want to:") {
		t.Fatalf("unexpected ambiguity prompt when active project is set: %s", output)
	}

	planView := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(ctx, "show the plan")
	})
	if !strings.Contains(strings.ToLower(planView), "planning") && !strings.Contains(strings.ToLower(planView), "tasks") {
		t.Fatalf("plan view should list stages/tasks, got: %s", planView)
	}
	if strings.Contains(planView, proj.ID) {
		t.Fatalf("plan output must not include uuid, got %s", planView)
	}

	details, err := svc.GetProjectDetails(ctx, proj.ID, user)
	if err != nil {
		t.Fatalf("project details: %v", err)
	}
	if len(details.Stages) == 0 || len(details.Tasks) == 0 {
		t.Fatalf("plan should create stages and tasks, got %d stages %d tasks", len(details.Stages), len(details.Tasks))
	}
}

func TestPlanRequiresActiveProject(t *testing.T) {
	ctx := context.Background()
	dbConn := connectTestDB(t, "file:plan-requires-active.db?_foreign_keys=1")
	user := models.User{ID: "user-no-active", Name: "No Active"}
	dispatcher, state, _ := newTestDispatcher(t, dbConn, user)

	state.ActiveProjectID = ""
	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(ctx, "create a plan for Yusuf")
	})

	if !strings.Contains(output, noActiveProjectMessage()) {
		t.Fatalf("expected no-active-project guidance, got: %s", output)
	}
}

func TestProjectDeleteRequiresConfirmation(t *testing.T) {
	ctx := context.Background()
	dbConn := connectTestDB(t, "file:project-delete.db?_foreign_keys=1")
	user := models.User{ID: "user-delete", Name: "Delete User"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, user)

	proj, _, _ := svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "minecraft speedrun"}, user)
	_, _ = dispatcher.Dispatch(ctx, "/project list")
	_, _ = dispatcher.Dispatch(ctx, "/project use 1")

	first := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(ctx, "/project delete 1")
	})
	if !strings.Contains(first, "Delete project") {
		t.Fatalf("expected confirmation prompt, got %s", first)
	}
	if state.PendingDeletion == nil {
		t.Fatalf("pending deletion should be set")
	}

	second := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(ctx, fmt.Sprintf("delete %s", proj.Title))
	})
	if strings.TrimSpace(state.ActiveProjectID) != "" {
		t.Fatalf("active project should be cleared after deletion")
	}
	if !strings.Contains(second, "deleted") {
		t.Fatalf("expected deletion acknowledgement, got %s", second)
	}

	projects, err := svc.ListProjects(ctx, user, services.ListFilters{})
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("deleted project should not be listed, got %d", len(projects))
	}
}

func TestProjectQAUsesDB(t *testing.T) {
	ctx := context.Background()
	dbConn := connectTestDB(t, "file:project-qa.db?_foreign_keys=1")
	user := models.User{ID: "user-qa", Name: "QA User"}
	dispatcher, _, svc := newTestDispatcher(t, dbConn, user)

	proj, _, _ := svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "minecraft speedrun"}, user)
	plans := []services.StagePlan{{Title: "Planning", Tasks: []services.TaskPlan{{Title: "Kickoff", Priority: "high"}, {Title: "Setup"}}}}
	_, _, _ = svc.CreateStagesAndTasksFromChat(ctx, proj.ID, plans, user)

	_, _ = dispatcher.Dispatch(ctx, "/project list")
	_, _ = dispatcher.Dispatch(ctx, "/project use 1")

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(ctx, "show project")
	})

	if strings.Contains(output, proj.ID) {
		t.Fatalf("project details should not include uuid, got %s", output)
	}
	if !strings.Contains(strings.ToLower(output), "planning") || !strings.Contains(strings.ToLower(output), "kickoff") {
		t.Fatalf("expected stages and tasks in output, got %s", output)
	}
}
