package chatbot

import "testing"

func TestDetectIntentCreateProject(t *testing.T) {
	match := DetectIntent("Create a project for a 5-story residential building")
	if match.Intent != IntentCreateProject {
		t.Fatalf("expected create project, got %s", match.Intent)
	}
	if match.ProjectTitle == "" {
		t.Fatalf("expected project title to be extracted")
	}
}

func TestDetectIntentAddStages(t *testing.T) {
	match := DetectIntent("Add stages: Design; Build | tasks: permits, drawings")
	if match.Intent != IntentAddStagesAndTasks {
		t.Fatalf("expected add stages intent, got %s", match.Intent)
	}
	if len(match.StagePlans) == 0 {
		t.Fatalf("expected stage plans")
	}
}

func TestDetectIntentAssign(t *testing.T) {
	match := DetectIntent("Assign John to task permits")
	if match.Intent != IntentAssignResponsible {
		t.Fatalf("expected assign intent")
	}
	if match.AssigneeName != "john" && match.AssigneeName == "" {
		t.Fatalf("assignee name not captured")
	}
}

func TestDetectIntentList(t *testing.T) {
	match := DetectIntent("list projects")
	if match.Intent != IntentListProjects {
		t.Fatalf("expected list intent")
	}
}

func TestDetectIntentListVariants(t *testing.T) {
	inputs := []string{"show projects", "what projects do I have?", "/project list"}
	for _, in := range inputs {
		match := DetectIntent(in)
		if match.Intent != IntentListProjects {
			t.Fatalf("input %q expected list intent, got %s", in, match.Intent)
		}
	}
}

func TestDetectIntentShow(t *testing.T) {
	match := DetectIntent("show project Skyline")
	if match.Intent != IntentShowProject {
		t.Fatalf("expected show intent")
	}
	if match.ProjectTitle == "" {
		t.Fatalf("expected project title")
	}
}

func TestDetectIntentShowPlan(t *testing.T) {
	match := DetectIntent("show the plan")
	if match.Intent != IntentShowProject {
		t.Fatalf("expected show project intent for plan view")
	}
}
