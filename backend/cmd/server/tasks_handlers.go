package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	dbpkg "tmplatform-backend/internal/db"

	"github.com/go-chi/chi/v5"
)

// Helper to get auth info from context (assuming authMiddleware sets it)
func (a *App) getAuthCtx(ctx context.Context) (*AuthCtx, bool) {
	v, ok := ctx.Value(authCtxKey).(*AuthCtx)
	return v, ok
}

// Handler: GET /stages/{id}/tasks
func (a *App) handleStageTasksList(w http.ResponseWriter, r *http.Request) {
	stageID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	tasks, err := dbpkg.ListTasksByStage(r.Context(), a.DB, stageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

// Handler: POST /tasks
func (a *App) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StageID     int    `json:"stage_id"`
		ParentID    *int   `json:"parent_id"` // Optional: for subtasks
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
		AssigneeID  *int   `json:"assignee_id"`
		DueDate     string `json:"due_date"` // YYYY-MM-DD
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	auth, ok := a.getAuthCtx(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse date if present
	var due *time.Time
	if input.DueDate != "" {
		t, err := time.Parse("2006-01-02", input.DueDate)
		if err == nil {
			due = &t
		}
	}

	var desc *string
	if input.Description != "" {
		desc = &input.Description
	}

	if input.AssigneeID == nil {
		input.AssigneeID = &auth.UserID
	}

	task := &dbpkg.Task{
		StageID:     input.StageID,
		ParentID:    input.ParentID,
		AuthorID:    auth.UserID,
		Title:       input.Title,
		Description: desc,
		Status:      "todo", // Default status
		Priority:    input.Priority,
		AssigneeID:  input.AssigneeID,
		DueDate:     due,
	}

	if err := dbpkg.CreateTask(r.Context(), a.DB, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

// Handler: GET /tasks?scope=my|urgent|subordinate
func (a *App) handleTasksList(w http.ResponseWriter, r *http.Request) {
	auth, ok := a.getAuthCtx(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	scope := r.URL.Query().Get("scope")
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if x, err := strconv.Atoi(v); err == nil {
			limit = x
		}
	}

	items, err := dbpkg.ListTasksForDashboard(r.Context(), a.DB, auth.OrgID, auth.UserID, scope, limit)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

// Handler: GET /tasks/{id}
func (a *App) handleTaskGet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	task, err := dbpkg.GetTaskByID(r.Context(), a.DB, id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// Fetch subtasks
	subtasks, _ := dbpkg.ListSubtasks(r.Context(), a.DB, id)

	// Fetch files
	files, _ := dbpkg.ListFilesByTask(r.Context(), a.DB, id)

	resp := map[string]interface{}{
		"task":     task,
		"subtasks": subtasks,
		"files":    files,
	}
	writeJSON(w, http.StatusOK, resp)
}

// Handler: POST /tasks/{id}/files
func (a *App) handleTaskFileUpload(w http.ResponseWriter, r *http.Request) {
	taskID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	// Max upload size 50MB
	r.ParseMultipartForm(50 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "bad file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	auth, ok := a.getAuthCtx(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Save file locally
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	// Unique filename: taskID_timestamp_originalName
	filename := fmt.Sprintf("%d_%d_%s", taskID, time.Now().Unix(), header.Filename)
	path := filepath.Join(uploadDir, filename)

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "write error", http.StatusInternalServerError)
		return
	}

	fObj := &dbpkg.File{
		TaskID:      taskID,
		UploaderID:  auth.UserID,
		Filename:    header.Filename,
		StoragePath: path,
		MimeType:    header.Header.Get("Content-Type"),
		SizeBytes:   header.Size,
	}

	if err := dbpkg.CreateFile(r.Context(), a.DB, fObj); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, fObj)
}
