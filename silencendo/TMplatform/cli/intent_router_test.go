package cli

import (
	"bytes"
	"context"
	stdcontext "context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"silencendo/chatbot"
	botcontext "silencendo/context"
	"silencendo/db"
	"silencendo/ingestion"
	"silencendo/llm"
	"silencendo/models"
	"silencendo/services"
	"silencendo/sources"
)

func TestProjectIntentDoesNotEditDocuments(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	filePath := filepath.Join(currentDir(t), "doc.txt")
	if err := os.WriteFile(filePath, []byte("original content"), 0644); err != nil {
		t.Fatalf("write doc: %v", err)
	}

	dbConn := connectTestDB(t, "file:project-intent.db?_foreign_keys=1")
	actor := models.User{ID: "user-project", Name: "Project User"}
	dispatcher, _, _ := newTestDispatcher(t, dbConn, actor)
	doc := dispatcher.documentHandler

	sm := sources.NewSourceManager()
	doc.sourceManager = sm
	doc.sourceCommands = sources.NewSourceCommands(sm)
	sm.AddSource("file", filePath)

	before := modTime(t, filePath)
	if _, err := dispatcher.Dispatch(context.Background(), "create project Demo Build"); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	after := modTime(t, filePath)

	if !after.Equal(before) {
		t.Fatalf("expected document untouched for project intent")
	}

	if count := countProjects(t, dbConn); count == 0 {
		t.Fatalf("expected a project to be created")
	}
}

func TestDocumentIntentDoesNotTouchDB(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	filePath := filepath.Join(currentDir(t), "doc.txt")
	if err := os.WriteFile(filePath, []byte("doc text to improve"), 0644); err != nil {
		t.Fatalf("write doc: %v", err)
	}

	dbConn := connectTestDB(t, "file:document-intent.db?_foreign_keys=1")
	actor := models.User{ID: "user-doc", Name: "Doc User"}
	dispatcher, _, _ := newTestDispatcher(t, dbConn, actor)
	dispatcher.documentHandler.llmClient = llm.NewMockLLM()

	sm := sources.NewSourceManager()
	dispatcher.documentHandler.sourceManager = sm
	dispatcher.documentHandler.sourceCommands = sources.NewSourceCommands(sm)
	src := sm.AddSource("file", filePath)
	sm.SetActiveSources([]string{src.ID})

	if _, err := dispatcher.Dispatch(context.Background(), "improve this document with better wording"); err != nil {
		t.Fatalf("dispatch: %v", err)
	}

	if count := countProjects(t, dbConn); count != 0 {
		t.Fatalf("expected no projects to be created during document edit, got %d", count)
	}
}

func TestProjectScopedWithoutActiveShowsGuidance(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:ambiguous-intent.db?_foreign_keys=1")
	actor := models.User{ID: "user-amb", Name: "Ambiguous User"}
	dispatcher, state, _ := newTestDispatcher(t, dbConn, actor)
	sm := sources.NewSourceManager()
	dispatcher.documentHandler.sourceManager = sm
	dispatcher.documentHandler.sourceCommands = sources.NewSourceCommands(sm)
	dispatcher.router.sourceManager = sm
	src := sm.AddSource("file", filepath.Join(currentDir(t), "doc.txt"))
	sm.SetActiveSources([]string{src.ID})

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "project plan")
	})

	if state.Pending != nil {
		t.Fatalf("pending confirmation should not be set for project-scoped input")
	}
	if !strings.Contains(output, noActiveProjectMessage()) {
		t.Fatalf("expected no-active-project guidance, got: %s", output)
	}
}

func TestProjectsPersistAcrossRestart(t *testing.T) {
	dir := t.TempDir()
	dsn := "file:" + filepath.Join(dir, "persistent.db") + "?_foreign_keys=1"

	conn1 := connectTestDB(t, dsn)
	svc1 := services.NewProjectService(conn1)
	user := models.User{ID: "persist-user", Name: "Persistent"}

	if _, created, err := svc1.CreateProjectFromChat(context.Background(), services.CreateProjectInput{Title: "Bridge"}, user); err != nil || !created {
		t.Fatalf("create project: %v", err)
	}
	conn1.Close()

	conn2 := connectTestDB(t, dsn)
	svc2 := services.NewProjectService(conn2)
	projects, err := svc2.ListProjects(context.Background(), user, services.ListFilters{})
	if err != nil {
		t.Fatalf("list projects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected persisted project count 1, got %d", len(projects))
	}
}

func TestStartHasNoModeSelectionPrompt(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:no-mode.db?_foreign_keys=1")
	actor := models.User{ID: "user-start", Name: "Starter"}
	dispatcher, _, _ := newTestDispatcher(t, dbConn, actor)
	repl := NewUnifiedRepl(dispatcher)

	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	reader := strings.NewReader("/exit\n")

	done := make(chan error, 1)
	go func() {
		done <- repl.Start(context.Background(), reader)
	}()

	wOut.Close()
	outBytes, _ := io.ReadAll(rOut)
	if err := <-done; err != nil {
		t.Fatalf("start error: %v", err)
	}

	output := strings.ToLower(string(outBytes))
	disallowed := []string{"main silencendo", "tmplatform", "choose", "select"}
	for _, word := range disallowed {
		if strings.Contains(output, word) {
			t.Fatalf("found disallowed startup prompt %q in output: %s", word, output)
		}
	}
}

func TestProjectQuestionsBypassLLMAndUseDB(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-llm-bypass.db?_foreign_keys=1")
	actor := models.User{ID: "user-proj-bypass", Name: "Proj Bypass"}
	dispatcher, _, svc := newTestDispatcher(t, dbConn, actor)

	// Create projects directly via service to avoid LLM/doc paths.
	ctx := context.Background()
	if _, _, err := svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "Alpha"}, actor); err != nil {
		t.Fatalf("create alpha: %v", err)
	}
	if _, _, err := svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "Beta"}, actor); err != nil {
		t.Fatalf("create beta: %v", err)
	}

	// Replace document handler LLM with a failing stub to ensure it is not called.
	failingLLM := &failLLM{t: t}
	dispatcher.documentHandler.llmClient = failingLLM

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(ctx, "what projects do we have?")
	})

	if strings.Contains(output, "LLM") {
		t.Fatalf("LLM placeholder should not appear for project queries: %s", output)
	}
	if !strings.Contains(output, "Alpha") || !strings.Contains(output, "Beta") {
		t.Fatalf("expected DB projects listed, got %s", output)
	}
}

func TestProjectQuestionsEmptyList(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-empty.db?_foreign_keys=1")
	actor := models.User{ID: "user-proj-empty", Name: "Proj Empty"}
	dispatcher, _, _ := newTestDispatcher(t, dbConn, actor)

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "what projects do we have?")
	})

	expected := "You have no projects yet. Create one with 'create project <name>'."
	if !strings.Contains(output, expected) {
		t.Fatalf("expected empty project message, got %s", output)
	}
}

func TestProjectListQueriesBypassLLM(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-list-queries.db?_foreign_keys=1")
	actor := models.User{ID: "user-list", Name: "Lister"}
	dispatcher, _, svc := newTestDispatcher(t, dbConn, actor)

	ctx := context.Background()
	if _, _, err := svc.CreateProjectFromChat(ctx, services.CreateProjectInput{Title: "minecraft speedrun"}, actor); err != nil {
		t.Fatalf("create project: %v", err)
	}

	failing := &failLLM{t: t}
	dispatcher.documentHandler.llmClient = failing

	inputs := []string{
		"what project do we have?",
		"do I have any projects?",
		"my projects",
	}

	for _, in := range inputs {
		output := captureOutput(t, func() {
			_, _ = dispatcher.Dispatch(ctx, in)
		})
		if !strings.Contains(output, "minecraft speedrun") {
			t.Fatalf("expected project list in output for %q, got: %s", in, output)
		}
		if strings.Contains(output, "LLM") {
			t.Fatalf("LLM output should not appear for %q, got: %s", in, output)
		}
	}
}

func TestDeepSeekClassificationCreatesPlanWithActiveProject(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-plan-deepseek.db?_foreign_keys=1")
	actor := models.User{ID: "user-plan", Name: "Planner"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, actor)

	project, _, err := svc.CreateProjectFromChat(context.Background(), services.CreateProjectInput{Title: "minecraft speedrun"}, actor)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	state.ActiveProjectID = project.ID
	state.ActiveProjectTitle = project.Title

	classifier := &fakeClassifierLLM{resp: `{"domain":"project","intent":"project_create_plan","confidence":0.9}`}
	dispatcher.router.llmClient = classifier
	dispatcher.router.useLLM = true

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "create a step-by-step plan for 2 players Yussuf and Omar to beat minecraft")
	})

	if !strings.Contains(strings.ToLower(output), "plan created") {
		t.Fatalf("expected plan creation, got: %s", output)
	}
	if classifier.answerCalls > 0 || classifier.editCalls > 0 {
		t.Fatalf("execution should not call Answer/Edit, got answer=%d edit=%d", classifier.answerCalls, classifier.editCalls)
	}
}

func TestDeepSeekClassificationShowsPlanWithActiveProject(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-plan-show.db?_foreign_keys=1")
	actor := models.User{ID: "user-plan-show", Name: "Planner"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, actor)

	project, _, err := svc.CreateProjectFromChat(context.Background(), services.CreateProjectInput{Title: "minecraft speedrun"}, actor)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	state.ActiveProjectID = project.ID
	state.ActiveProjectTitle = project.Title

	// Seed a plan so show_plan has content.
	_, _, _, err = svc.CreateRuleBasedPlan(context.Background(), project.ID, nil, actor)
	if err != nil {
		t.Fatalf("seed plan: %v", err)
	}

	classifier := &fakeClassifierLLM{resp: `{"domain":"project","intent":"project_show_plan","confidence":0.92}`}
	dispatcher.router.llmClient = classifier
	dispatcher.router.useLLM = true

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "show me the content of this project")
	})

	if !strings.Contains(output, "Tasks:") {
		t.Fatalf("expected tasks to be shown, got: %s", output)
	}
}

func TestCreatePlanWithoutActiveProjectRequestsSelection(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-plan-no-active.db?_foreign_keys=1")
	actor := models.User{ID: "user-no-active", Name: "NoActive"}
	dispatcher, _, _ := newTestDispatcher(t, dbConn, actor)

	classifier := &fakeClassifierLLM{resp: `{"domain":"project","intent":"project_create_plan","confidence":0.88}`}
	dispatcher.router.llmClient = classifier
	dispatcher.router.useLLM = true

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "create a step by step plan for us")
	})

	if !strings.Contains(output, noActiveProjectMessage()) {
		t.Fatalf("expected guidance to select project, got: %s", output)
	}
}

func TestDeepSeekSuccessBeatsKeywordFallback(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-deepseek-priority.db?_foreign_keys=1")
	actor := models.User{ID: "user-priority", Name: "Priority"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, actor)

	project, _, err := svc.CreateProjectFromChat(context.Background(), services.CreateProjectInput{Title: "minecraft speedrun"}, actor)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	state.ActiveProjectID = project.ID
	state.ActiveProjectTitle = project.Title

	classifier := &fakeClassifierLLM{resp: `{"domain":"project","intent":"project_create_plan","confidence":0.95}`}
	dispatcher.router.llmClient = classifier
	dispatcher.router.useLLM = true

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "just wondering, maybe tell me something?")
	})

	if !strings.Contains(strings.ToLower(output), "plan created") {
		t.Fatalf("expected DeepSeek classification to trigger plan creation despite weak keywords, got: %s", output)
	}
}

func TestDeepSeekFailureDoesNotFallbackToKeywords(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-fallback.db?_foreign_keys=1")
	actor := models.User{ID: "user-fallback", Name: "Fallback"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, actor)

	project, _, err := svc.CreateProjectFromChat(context.Background(), services.CreateProjectInput{Title: "minecraft speedrun"}, actor)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	state.ActiveProjectID = project.ID
	state.ActiveProjectTitle = project.Title

	classifier := &fakeClassifierLLM{resp: "<<<not json>>>"}
	dispatcher.router.llmClient = classifier
	dispatcher.router.useLLM = true

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "create plan for minecraft")
	})

	if !strings.Contains(strings.ToLower(output), "deepseek couldn't understand") {
		t.Fatalf("expected deepseek failure notice instead of keyword fallback, got: %s", output)
	}
}

func TestProjectExecutionDoesNotCallLLMDuringHandlers(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:project-no-llm-exec.db?_foreign_keys=1")
	actor := models.User{ID: "user-no-llm", Name: "NoLLM"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, actor)

	project, _, err := svc.CreateProjectFromChat(context.Background(), services.CreateProjectInput{Title: "minecraft speedrun"}, actor)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	state.ActiveProjectID = project.ID
	state.ActiveProjectTitle = project.Title

	classifier := &fakeClassifierLLM{resp: `{"domain":"project","intent":"project_create_plan","confidence":0.9}`}
	dispatcher.router.llmClient = classifier
	dispatcher.router.useLLM = true

	_ = captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "create plan now")
	})

	if classifier.answerCalls != 0 || classifier.editCalls != 0 {
		t.Fatalf("Answer/Edit should not be called during project execution, got answer=%d edit=%d", classifier.answerCalls, classifier.editCalls)
	}
	if classifier.generateCalls == 0 {
		t.Fatalf("expected Generate to be called for classification")
	}
}

func TestResolvedProjectSkipsAmbiguity(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:resolved-project.db?_foreign_keys=1")
	actor := models.User{ID: "user-resolved", Name: "Resolver"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, actor)
	project, _, err := svc.CreateProjectFromChat(context.Background(), services.CreateProjectInput{Title: "minecraft speedrun"}, actor)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	state.ActiveProjectID = project.ID
	state.ActiveProjectTitle = project.Title

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "Create a plan for minecraft speedrun")
	})

	if strings.Contains(output, "Do you want to:") {
		t.Fatalf("expected no ambiguity prompt when project resolved, got: %s", output)
	}
	if state.Pending != nil {
		t.Fatalf("pending confirmation should not be set when project resolved")
	}
}

func TestAmbiguityWhenNoProjectReference(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:resolved-ambiguous.db?_foreign_keys=1")
	actor := models.User{ID: "user-ambig", Name: "Ambiguous"}
	dispatcher, state, _ := newTestDispatcher(t, dbConn, actor)

	sm := sources.NewSourceManager()
	dispatcher.documentHandler.sourceManager = sm
	dispatcher.documentHandler.sourceCommands = sources.NewSourceCommands(sm)
	dispatcher.router.sourceManager = sm
	src := sm.AddSource("file", filepath.Join(currentDir(t), "doc.txt"))
	sm.SetActiveSources([]string{src.ID})

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "Create a plan for a game")
	})

	if !strings.Contains(output, noActiveProjectMessage()) {
		t.Fatalf("expected no active project message, got: %s", output)
	}
	if state.Pending != nil {
		t.Fatalf("pending confirmation should remain nil when guidance is shown")
	}
}

func TestMixedProjectVerbsStayProject(t *testing.T) {
	restore := useTempWorkdir(t)
	defer restore()

	dbConn := connectTestDB(t, "file:mixed-project.db?_foreign_keys=1")
	actor := models.User{ID: "user-mixed", Name: "Mixed"}
	dispatcher, state, svc := newTestDispatcher(t, dbConn, actor)
	project, _, err := svc.CreateProjectFromChat(context.Background(), services.CreateProjectInput{Title: "minecraft speedrun"}, actor)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	state.ActiveProjectID = project.ID
	state.ActiveProjectTitle = project.Title

	output := captureOutput(t, func() {
		_, _ = dispatcher.Dispatch(context.Background(), "Add people Yusuf and Omar to minecraft speedrun and create plan")
	})

	if strings.Contains(output, "Do you want to:") {
		t.Fatalf("expected no ambiguity prompt for mixed verbs with project, got: %s", output)
	}
	if state.Pending != nil {
		t.Fatalf("pending confirmation should not be set for resolved project input")
	}
	if strings.Contains(strings.ToLower(output), "edit result") {
		t.Fatalf("LLM edit output should not appear for project intent: %s", output)
	}
}

type failLLM struct {
	t *testing.T
}

func (f *failLLM) Generate(stdcontext.Context, []llm.Message) (string, error) {
	f.t.Fatalf("LLM Generate should not be called for project queries")
	return "", fmt.Errorf("should not be called")
}

func (f *failLLM) Answer(stdcontext.Context, string, []ingestion.Chunk, botcontext.ContextType) (string, error) {
	f.t.Fatalf("LLM Answer should not be called for project queries")
	return "", fmt.Errorf("should not be called")
}

func (f *failLLM) Edit(stdcontext.Context, string, string, []ingestion.Chunk, botcontext.ContextType) (string, error) {
	f.t.Fatalf("LLM Edit should not be called for project queries")
	return "", fmt.Errorf("should not be called")
}

type fakeClassifierLLM struct {
	resp          string
	generateCalls int
	answerCalls   int
	editCalls     int
}

func (f *fakeClassifierLLM) Generate(stdcontext.Context, []llm.Message) (string, error) {
	f.generateCalls++
	return f.resp, nil
}

func (f *fakeClassifierLLM) Answer(stdcontext.Context, string, []ingestion.Chunk, botcontext.ContextType) (string, error) {
	f.answerCalls++
	return "", nil
}

func (f *fakeClassifierLLM) Edit(stdcontext.Context, string, string, []ingestion.Chunk, botcontext.ContextType) (string, error) {
	f.editCalls++
	return "", nil
}

func newTestDispatcher(t *testing.T, dbConn *sql.DB, actor models.User) (*Dispatcher, *SessionState, *services.ProjectService) {
	t.Helper()
	state := NewSessionState()
	docHandler := NewCLIInterface()
	projectService := services.NewProjectService(dbConn)
	projectHandler := chatbot.NewHandler(projectService, actor)
	router := NewIntentRouter(docHandler.SourceManager(), docHandler.IntentDetector())
	dispatcher := NewDispatcher(router, projectHandler, docHandler, state)
	return dispatcher, state, projectService
}

func connectTestDB(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	cfg := db.Config{Driver: "sqlite", DSN: dsn}
	conn, err := db.Connect(cfg)
	if err != nil {
		t.Fatalf("connect db: %v", err)
	}
	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	return conn
}

func countProjects(t *testing.T, dbConn *sql.DB) int {
	t.Helper()
	var count int
	if err := dbConn.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count); err != nil {
		t.Fatalf("count projects: %v", err)
	}
	return count
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	outBytes, _ := io.ReadAll(r)
	os.Stdout = old
	return string(outBytes)
}

func useTempWorkdir(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	return func() {
		_ = os.Chdir(cwd)
	}
}

func currentDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return dir
}

func modTime(t *testing.T, path string) time.Time {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return info.ModTime()
}

// Interface guard for unused bytes import in this file.
var _ = bytes.NewBuffer
