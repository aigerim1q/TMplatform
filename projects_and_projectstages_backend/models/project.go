package models

type Stage struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	ProjectID int      `json:"project_id"`
	CreatedAt string   `json:"created_at"`
	Status    string   `json:"status"`
	Budget    float64  `json:"budget"`
	TasksIDs  []int    `json:"tasks_ids"`
	Project   *Project `json:"project,omitempty"`
}

type Project struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	LifecycleID   *int    `json:"lifecycle_id"`
	CreatedAt     string  `json:"created_at"`
	Status        string  `json:"status"`
	ResponsibleID int     `json:"responsible_id"`
	Budget        float64 `json:"budget"`
	Stages        []Stage `json:"stages"`
}

type Task struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	StageID int    `json:"stage_id"`
}

type Subtask struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	TaskID int    `json:"task_id"`
}
