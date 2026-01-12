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

func TestDetectIntentShowContent(t *testing.T) {
	match := DetectIntent("show the content")
	if match.Intent != IntentShowProject {
		t.Fatalf("expected show project intent for content")
	}
}

func TestDetectIntentWhatIsContent(t *testing.T) {
	match := DetectIntent("what is the content")
	if match.Intent != IntentShowProject {
		t.Fatalf("expected show project intent for 'what is the content'")
	}
}

func TestExtractPlanMembersSplitsAnd(t *testing.T) {
	members := extractPlanMembers("create plan for Omar and Yussuf to beat minecraft")
	if len(members) != 2 {
		t.Fatalf("expected two members, got %d", len(members))
	}
	if members[0] != "Omar" || members[1] != "Yussuf" {
		t.Fatalf("unexpected members parsed: %#v", members)
	}
}

func TestExtractPlanMembersWhitespaceSeparated(t *testing.T) {
	members := extractPlanMembers("create plan for Omar Yussuf Fatima to beat minecraft")
	if len(members) != 3 {
		t.Fatalf("expected three members, got %d", len(members))
	}
	if members[0] != "Omar" || members[1] != "Yussuf" || members[2] != "Fatima" {
		t.Fatalf("unexpected members parsed: %#v", members)
	}
}
