package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"backend-go/db"
	"backend-go/models"

	"github.com/gorilla/mux"
)

var PROJECT_STATUSES = []string{"не начат", "в работе", "на паузе", "завершен"}

func getProjectByIDDB(id int) (*models.Project, error) {
	project := &models.Project{}

	err := db.DB.QueryRow(`
        SELECT id, name, lifecycle_id, created_at, status, responsible_id, budget
        FROM projects WHERE id=$1
    `, id).Scan(
		&project.ID,
		&project.Name,
		&project.LifecycleID,
		&project.CreatedAt,
		&project.Status,
		&project.ResponsibleID,
		&project.Budget,
	)
	if err != nil {
		return nil, err
	}

	stageRows, err := db.DB.Query(`
        SELECT id, project_id, name, created_at, status, budget
        FROM stages WHERE project_id=$1 ORDER BY id
    `, project.ID)
	if err != nil {
		return nil, err
	}
	defer stageRows.Close()

	stages := []models.Stage{}
	for stageRows.Next() {
		var stage models.Stage
		err := stageRows.Scan(
			&stage.ID,
			&stage.ProjectID,
			&stage.Name,
			&stage.CreatedAt,
			&stage.Status,
			&stage.Budget,
		)
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
		stages = append(stages, stage)
	}

	project.Stages = stages
	return project, nil
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		UserID        int     `json:"user_id"`
		Name          string  `json:"name"`
		ResponsibleID int     `json:"responsible_id"`
		Budget        float64 `json:"budget"`
	}

	var body ReqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var roleID int
	err := db.DB.QueryRow("SELECT role_id FROM users WHERE id=$1", body.UserID).Scan(&roleID)
	if err != nil || (roleID != 1 && roleID != 2) {
		http.Error(w, "Only admins or managers can create projects", http.StatusForbidden)
		return
	}

	row := db.DB.QueryRow(`
        INSERT INTO projects (name, responsible_id, budget, status)
        VALUES ($1, $2, $3, 'не начат')
        RETURNING id, name, lifecycle_id, created_at, status, responsible_id, budget
    `, body.Name, body.ResponsibleID, body.Budget)

	project := models.Project{}
	err = row.Scan(
		&project.ID, &project.Name, &project.LifecycleID,
		&project.CreatedAt, &project.Status, &project.ResponsibleID, &project.Budget,
	)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	stageRow := db.DB.QueryRow(`
        INSERT INTO stages (project_id, name, status, budget)
        VALUES ($1, 'Планирование', 'не начат', $2)
        RETURNING id, project_id, name, created_at, status, budget
    `, project.ID, project.Budget)

	stage := models.Stage{}
	stageRow.Scan(&stage.ID, &stage.ProjectID, &stage.Name, &stage.CreatedAt, &stage.Status, &stage.Budget)
	project.Stages = []models.Stage{stage}

	json.NewEncoder(w).Encode(project)
}

func GetMyProjects(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	rows, err := db.DB.Query(`
        SELECT id, name, lifecycle_id, created_at, status, responsible_id, budget
        FROM projects WHERE responsible_id=$1 ORDER BY created_at DESC
    `, userID)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	projects := []models.Project{}
	for rows.Next() {
		var p models.Project
		rows.Scan(&p.ID, &p.Name, &p.LifecycleID, &p.CreatedAt, &p.Status, &p.ResponsibleID, &p.Budget)
		projects = append(projects, p)
	}
	json.NewEncoder(w).Encode(projects)
}

func GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	project, err := getProjectByIDDB(id)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(project)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	type ReqBody struct {
		Name          *string  `json:"name"`
		Budget        *float64 `json:"budget"`
		Status        *string  `json:"status"`
		ResponsibleID *int     `json:"responsible_id"`
	}

	var body ReqBody
	json.NewDecoder(r.Body).Decode(&body)

	set := []string{}
	args := []interface{}{}
	i := 1

	if body.Name != nil {
		set = append(set, fmt.Sprintf("name=$%d", i))
		args = append(args, *body.Name)
		i++
	}
	if body.Budget != nil {
		set = append(set, fmt.Sprintf("budget=$%d", i))
		args = append(args, *body.Budget)
		i++
	}
	if body.Status != nil {
		valid := false
		for _, s := range PROJECT_STATUSES {
			if s == *body.Status {
				valid = true
				break
			}
		}
		if !valid {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}
		set = append(set, fmt.Sprintf("status=$%d", i))
		args = append(args, *body.Status)
		i++
	}
	if body.ResponsibleID != nil {
		set = append(set, fmt.Sprintf("responsible_id=$%d", i))
		args = append(args, *body.ResponsibleID)
		i++
	}
	if len(set) == 0 {
		http.Error(w, "Nothing to update", http.StatusBadRequest)
		return
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE projects SET %s WHERE id=$%d RETURNING id,name,lifecycle_id,created_at,status,responsible_id,budget`, join(set, ", "), i)
	row := db.DB.QueryRow(query, args...)

	project := models.Project{}
	err := row.Scan(&project.ID, &project.Name, &project.LifecycleID, &project.CreatedAt, &project.Status, &project.ResponsibleID, &project.Budget)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(project)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	row := db.DB.QueryRow(`DELETE FROM projects WHERE id=$1 RETURNING id,name,lifecycle_id,created_at,status,responsible_id,budget`, id)

	project := models.Project{}
	err := row.Scan(&project.ID, &project.Name, &project.LifecycleID, &project.CreatedAt, &project.Status, &project.ResponsibleID, &project.Budget)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Project deleted",
		"project": project,
	})
}

func join(arr []string, sep string) string {
	out := ""
	for i, v := range arr {
		if i > 0 {
			out += sep
		}
		out += v
	}
	return out
}
