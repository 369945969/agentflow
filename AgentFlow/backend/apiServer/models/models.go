package models

import "time"

type Skill struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	Description string    `json:"description" db:"description"`
	Version     string    `json:"version" db:"version"`
	Icon        string    `json:"icon" db:"icon"`
	Color       string    `json:"color" db:"color"`
	Logic       string    `json:"logic" db:"logic"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Group struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Members          []string  `json:"members"`
	GroupRuleMode    string    `json:"group_rule_mode,omitempty" db:"group_rule_mode"`
	ThinkingEnabled  bool      `json:"thinking_enabled,omitempty" db:"thinking_enabled"`
	SimplifiedOutput bool      `json:"simplified_output,omitempty" db:"simplified_output"`
	CustomRule       string    `json:"custom_rule,omitempty" db:"custom_rule"`
}

type Agent struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	Model            string    `json:"model" db:"model"`
	SystemPrompt     string    `json:"system_prompt" db:"system_prompt"`
	ThinkingEnabled  bool      `json:"thinking_enabled,omitempty" db:"thinking_enabled"`
	SimplifiedOutput bool      `json:"simplified_output,omitempty" db:"simplified_output"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Skills           []string  `json:"skills"` // Bound skill IDs
}

type Edge struct {
	ID        string    `json:"id" db:"id"`
	FromID    string    `json:"from_id" db:"from_id"`
	ToID      string    `json:"to_id" db:"to_id"`
	Type      string    `json:"type" db:"type"`
	Metadata  string    `json:"metadata" db:"metadata"` // JSON as string in DuckDB
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Model struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Provider    string    `json:"provider" db:"provider"`
	BaseURL     string    `json:"base_url" db:"base_url"`
	APIKey      string    `json:"api_key" db:"api_key"`
	ModelName   string    `json:"model_name" db:"model_name"`
	Description string    `json:"description" db:"description"`
	IsDefault   bool      `json:"is_default" db:"is_default"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
