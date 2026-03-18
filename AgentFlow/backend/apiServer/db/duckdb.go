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
		dbPath = "./data/agentflow.duckdb"
	}
	
	// Create data directory if it doesn't exist
	dataDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	var err error
	DB, err = sql.Open("duckdb", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to DuckDB: %v", err)
	}

	// Initialize tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS skills (id VARCHAR PRIMARY KEY, name VARCHAR, type VARCHAR, description TEXT, version VARCHAR, icon VARCHAR, color VARCHAR, logic TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS agents (id VARCHAR PRIMARY KEY, name VARCHAR, description TEXT, model VARCHAR, system_prompt TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS edges (id VARCHAR PRIMARY KEY, from_id VARCHAR, to_id VARCHAR, type VARCHAR, metadata JSON, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS models (id VARCHAR PRIMARY KEY, name VARCHAR, provider VARCHAR, base_url VARCHAR, api_key VARCHAR, model_name VARCHAR, description TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatalf("Failed to initialize database table: %v", err)
		}
	}

	fmt.Println("✅ DuckDB Initialized (Relational + Graph)")
}
