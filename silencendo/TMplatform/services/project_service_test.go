package services

import (
	"context"
	"testing"

	"silencendo/models"
)

func TestCreateProjectIdempotent(t *testing.T) {
	svc := newTestService(t)
	user := models.User{ID: "user-1", Name: "Alice"}

	ctx := context.Background()

	p1, created, err := svc.CreateProjectFromChat(ctx, CreateProjectInput{Title: "Skyline"}, user)
	if err != nil {
		t.Fatalf("create project error: %v", err)
	}
	if !created {
		t.Fatalf("expected creation")
	}

	p2, created2, err := svc.CreateProjectFromChat(ctx, CreateProjectInput{Title: "Skyline"}, user)
	if err != nil {
		t.Fatalf("idempotent create error: %v", err)
	}
	if created2 {
		t.Fatalf("expected second call to be idempotent")
	}
	if p1.ID != p2.ID {
		t.Fatalf("expected same project id, got %s vs %s", p1.ID, p2.ID)
	}
}

func TestAddStagesAndTasks(t *testing.T) {
	svc := newTestService(t)
	user := models.User{ID: "user-2", Name: "Bob"}
	ctx := context.Background()

	project, created, err := svc.CreateProjectFromChat(ctx, CreateProjectInput{Title: "Bridge"}, user)
	if err != nil || !created {
		t.Fatalf("project create failed: %v", err)
	}

	plans := []StagePlan{{
		Title: "Design",
		Tasks: []TaskPlan{{Title: "Permits"}, {Title: "Drawings"}},
	}}

	stages, tasks, err := svc.CreateStagesAndTasksFromChat(ctx, project.ID, plans, user)
	if err != nil {
		t.Fatalf("add stages/tasks failed: %v", err)
	}
	if len(stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(stages))
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	// Idempotency: re-run
	_, tasks2, err := svc.CreateStagesAndTasksFromChat(ctx, project.ID, plans, user)
	if err != nil {
		t.Fatalf("idempotent stage/task call failed: %v", err)
	}
	if len(tasks2) != 2 {
		t.Fatalf("expected idempotent tasks to remain 2, got %d", len(tasks2))
	}
}

func TestListProjectsAccessControl(t *testing.T) {
	svc := newTestService(t)
	owner := models.User{ID: "owner", Name: "Owner"}
	other := models.User{ID: "other", Name: "Other"}
	ctx := context.Background()

	project, _, err := svc.CreateProjectFromChat(ctx, CreateProjectInput{Title: "Private"}, owner)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	projects, err := svc.ListProjects(ctx, owner, ListFilters{})
	if err != nil {
		t.Fatalf("list owner projects: %v", err)
	}
	if len(projects) != 1 || projects[0].ID != project.ID {
		t.Fatalf("owner should see the project")
	}

	projectsOther, err := svc.ListProjects(ctx, other, ListFilters{})
	if err != nil {
		t.Fatalf("list other projects: %v", err)
	}
	if len(projectsOther) != 0 {
		t.Fatalf("non-member should not see projects")
	}
}

func TestListProjectsMultiple(t *testing.T) {
	svc := newTestService(t)
	user := models.User{ID: "user-multi", Name: "User Multi"}
	ctx := context.Background()

	if _, _, err := svc.CreateProjectFromChat(ctx, CreateProjectInput{Title: "Alpha"}, user); err != nil {
		t.Fatalf("create alpha: %v", err)
	}
	if _, _, err := svc.CreateProjectFromChat(ctx, CreateProjectInput{Title: "Beta"}, user); err != nil {
		t.Fatalf("create beta: %v", err)
	}

	projects, err := svc.ListProjects(ctx, user, ListFilters{})
	if err != nil {
		t.Fatalf("list projects: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}

	// Ensure both titles appear; ordering should be newest first.
	if projects[0].Title == projects[1].Title {
		t.Fatalf("expected distinct projects")
	}
	titles := []string{projects[0].Title, projects[1].Title}
	seenAlpha := titles[0] == "Alpha" || titles[1] == "Alpha"
	seenBeta := titles[0] == "Beta" || titles[1] == "Beta"
	if !seenAlpha || !seenBeta {
		t.Fatalf("expected both Alpha and Beta in list, got %v", titles)
	}
}

func TestListProjectsEmpty(t *testing.T) {
	svc := newTestService(t)
	user := models.User{ID: "user-empty", Name: "User Empty"}
	ctx := context.Background()

	projects, err := svc.ListProjects(ctx, user, ListFilters{})
	if err != nil {
		t.Fatalf("list projects: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("expected empty projects list, got %d", len(projects))
	}
}

func TestAssignResponsibleStage(t *testing.T) {
	svc := newTestService(t)
	owner := models.User{ID: "owner-assign", Name: "Owner"}
	ctx := context.Background()

	project, _, err := svc.CreateProjectFromChat(ctx, CreateProjectInput{Title: "AssignProj"}, owner)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	stages, _, err := svc.CreateStagesAndTasksFromChat(ctx, project.ID, []StagePlan{{Title: "Design"}}, owner)
	if err != nil || len(stages) == 0 {
		t.Fatalf("create stage: %v", err)
	}

	assignee := models.User{ID: "user-assignee", Name: "Assignee"}
	if _, err := svc.EnsureUser(ctx, assignee); err != nil {
		t.Fatalf("ensure assignee: %v", err)
	}

	if err := svc.AssignResponsible(ctx, "stage", stages[0].ID, assignee.ID, owner); err != nil {
		t.Fatalf("assign responsible failed: %v", err)
	}
}
