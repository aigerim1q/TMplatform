package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend-go/db"
	"backend-go/models"

	"github.com/gorilla/mux"
)

var STAGE_STATUSES = []string{"не начат", "в работе", "завершен"}

func getStageByIDDB(id int) (*models.Stage, error) {
	stage := &models.Stage{}
	err := db.DB.QueryRow(`
        SELECT id, project_id, name, created_at, status, budget
        FROM stages WHERE id=$1
    `, id).Scan(&stage.ID, &stage.ProjectID, &stage.Name, &stage.CreatedAt, &stage.Status, &stage.Budget)
	if err != nil {
		return nil, err
	}

	taskRows, _ := db.DB.Query(`SELECT id FROM tasks WHERE stage_id=$1`, stage.ID)
	taskIDs := []int{}
	for taskRows.Next() {
		var tid int
		taskRows.Scan(&tid)
		taskIDs = append(taskIDs, tid)
	}
	taskRows.Close()
	stage.TasksIDs = taskIDs

	return stage, nil
}

func CreateStage(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		ProjectID int     `json:"project_id"`
		Name      string  `json:"name"`
		Budget    float64 `json:"budget"`
	}
	var body ReqBody
	json.NewDecoder(r.Body).Decode(&body)

	row := db.DB.QueryRow(`
        INSERT INTO stages (project_id, name, status, budget)
        VALUES ($1, $2, 'не начат', $3)
        RETURNING id, project_id, name, created_at, status, budget
    `, body.ProjectID, body.Name, body.Budget)

	stage := models.Stage{}
	row.Scan(&stage.ID, &stage.ProjectID, &stage.Name, &stage.CreatedAt, &stage.Status, &stage.Budget)

	json.NewEncoder(w).Encode(stage)
}

func GetStagesByProject(w http.ResponseWriter, r *http.Request) {
	projectID, _ := strconv.Atoi(r.URL.Query().Get("project_id"))
	rows, _ := db.DB.Query(`
        SELECT id, project_id, name, created_at, status, budget
        FROM stages WHERE project_id=$1 ORDER BY id
    `, projectID)
	defer rows.Close()

	stages := []models.Stage{}
	for rows.Next() {
		var s models.Stage
		rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.CreatedAt, &s.Status, &s.Budget)

		taskRows, _ := db.DB.Query(`SELECT id FROM tasks WHERE stage_id=$1`, s.ID)
		taskIDs := []int{}
		for taskRows.Next() {
			var tid int
			taskRows.Scan(&tid)
			taskIDs = append(taskIDs, tid)
		}
		taskRows.Close()
		s.TasksIDs = taskIDs

		stages = append(stages, s)
	}

	json.NewEncoder(w).Encode(stages)
}

func GetStageByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	stage, err := getStageByIDDB(id)
	if err != nil {
		http.Error(w, "Stage not found", http.StatusNotFound)
		return
	}

	project, _ := getProjectByIDDB(stage.ProjectID)
	stage.Project = project

	json.NewEncoder(w).Encode(stage)
}

func UpdateStage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	type ReqBody struct {
		Status *string `json:"status"`
	}
	var body ReqBody
	json.NewDecoder(r.Body).Decode(&body)

	if body.Status == nil {
		http.Error(w, "Nothing to update", http.StatusBadRequest)
		return
	}

	valid := false
	for _, s := range STAGE_STATUSES {
		if s == *body.Status {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow(`
        UPDATE stages SET status=$1 WHERE id=$2
        RETURNING id, project_id, name, created_at, status, budget
    `, *body.Status, id)
	stage := models.Stage{}
	err := row.Scan(&stage.ID, &stage.ProjectID, &stage.Name, &stage.CreatedAt, &stage.Status, &stage.Budget)
	if err != nil {
		http.Error(w, "Stage not found", http.StatusNotFound)
		return
	}

	if *body.Status == "завершен" {
		nextRow := db.DB.QueryRow(`
            SELECT id FROM stages
            WHERE project_id=$1 AND id > $2
            ORDER BY id ASC LIMIT 1
        `, stage.ProjectID, stage.ID)
		var nextID int
		if nextRow.Scan(&nextID) == nil {
			db.DB.Exec("UPDATE stages SET status='не начат' WHERE id=$1", nextID)
		} else {
			db.DB.Exec("UPDATE projects SET status='завершен' WHERE id=$1", stage.ProjectID)
		}
	}

	if *body.Status == "в работе" {
		db.DB.Exec("UPDATE projects SET status='в работе' WHERE id=$1 AND status='не начат'", stage.ProjectID)
	}

	json.NewEncoder(w).Encode(stage)
}

func DeleteStage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	row := db.DB.QueryRow(`
        DELETE FROM stages WHERE id=$1
        RETURNING id, project_id, name, created_at, status, budget
    `, id)
	stage := models.Stage{}
	err := row.Scan(&stage.ID, &stage.ProjectID, &stage.Name, &stage.CreatedAt, &stage.Status, &stage.Budget)
	if err != nil {
		http.Error(w, "Stage not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Stage deleted",
		"stage":   stage,
	})
}
