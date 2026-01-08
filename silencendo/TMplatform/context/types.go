package context

type ContextType struct {
	Project    *string  `json:"project"`
	Stage      *string  `json:"stage"`
	Task       *string  `json:"task"`
	Goals      []string `json:"goals"`
	Constraints []string `json:"constraints"`
	Decisions   []string `json:"decisions"`
	NextSteps   []string `json:"next_steps"`
}

type ContextState struct {
	Active ContextType   `json:"active"`
	History []ContextType `json:"history"`
}