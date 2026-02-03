package prompt_engineering

// PromptTemplate represents a prompt template
type PromptTemplate struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Template    string   `json:"template"`
	Parameters  []string `json:"parameters"`
}

// PromptData holds data for prompt creation
type PromptData struct {
	DocumentContent string                 `json:"document_content"`
	JSONSchema      map[string]interface{} `json:"json_schema"`
}
