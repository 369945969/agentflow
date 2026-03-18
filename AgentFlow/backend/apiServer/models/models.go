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

type Agent struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Model        string    `json:"model" db:"model"`
	SystemPrompt string    `json:"system_prompt" db:"system_prompt"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
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
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

