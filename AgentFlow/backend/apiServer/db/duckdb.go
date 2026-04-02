package db

import (
	"apiServer/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/marcboeker/go-duckdb"
)

var DB *sql.DB

func InitDB() {
	dbPath := config.AppConfig.Database.Path
	if dbPath == "" {
		dbPath = "data/agentflow.duckdb"
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	var err error
	DB, err = sql.Open("duckdb", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Execute initial schema
	queries := []string{
		`CREATE TABLE IF NOT EXISTS skills (id VARCHAR PRIMARY KEY, name VARCHAR, type VARCHAR, description TEXT, version VARCHAR, icon VARCHAR, color VARCHAR, logic TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS models (id VARCHAR PRIMARY KEY, name VARCHAR, provider VARCHAR, base_url VARCHAR, api_key VARCHAR, model_name VARCHAR, description TEXT, is_default BOOLEAN DEFAULT FALSE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS agents (id VARCHAR PRIMARY KEY, name VARCHAR, description TEXT, model VARCHAR, system_prompt TEXT, thinking_enabled BOOLEAN DEFAULT TRUE, simplified_output BOOLEAN DEFAULT FALSE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS edges (id VARCHAR PRIMARY KEY, from_id VARCHAR, to_id VARCHAR, type VARCHAR, metadata JSON, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS groups (id VARCHAR PRIMARY KEY, name VARCHAR, group_rule_mode VARCHAR DEFAULT 'free', thinking_enabled BOOLEAN DEFAULT TRUE, simplified_output BOOLEAN DEFAULT FALSE, custom_rule TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS group_members (group_id VARCHAR NOT NULL, user_id VARCHAR NOT NULL, PRIMARY KEY (group_id, user_id));`,
		`CREATE TABLE IF NOT EXISTS agent_skills (agent_id VARCHAR NOT NULL, skill_id VARCHAR NOT NULL, PRIMARY KEY (agent_id, skill_id));`,
		`CREATE TABLE IF NOT EXISTS workflows (id VARCHAR PRIMARY KEY, name VARCHAR, description TEXT, graph JSON, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
	}
	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatalf("Failed to initialize database table: %v", err)
		}
	}

	fmt.Println("✅ DuckDB Initialized")
}
