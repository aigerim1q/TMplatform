package models

type StageModel struct {
	ID          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	StageTypeID *int   `json:"stage_type_id"`
	CreatedAt   string `json:"created_at"`
	Status      string `json:"status"`
	Tasks       []Task `json:"tasks"`
}
