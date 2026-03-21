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
		`CREATE TABLE IF NOT EXISTS models (id VARCHAR PRIMARY KEY, name VARCHAR, provider VARCHAR, base_url VARCHAR, api_key VARCHAR, model_name VARCHAR, description TEXT, is_default BOOLEAN DEFAULT FALSE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS groups (id VARCHAR PRIMARY KEY, name VARCHAR, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS group_members (group_id VARCHAR, user_id VARCHAR, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (group_id, user_id));`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatalf("Failed to initialize database table: %v", err)
		}
	}

	// Migrations: check if is_default column exists in models table
	_, err = DB.Exec("SELECT is_default FROM models LIMIT 1")
	if err != nil {
		log.Println("Migration: Adding is_default column to models table...")
		_, err = DB.Exec("ALTER TABLE models ADD COLUMN is_default BOOLEAN DEFAULT FALSE")
		if err != nil {
			log.Printf("Failed to migrate models table: %v", err)
		} else {
			log.Println("✅ Migration successful: is_default column added")
		}
	}

	// Migrations for groups table
	migrateGroupTable := func(columnName string, columnType string, defaultValue string) {
		_, err := DB.Exec(fmt.Sprintf("SELECT %s FROM groups LIMIT 1", columnName))
		if err != nil {
			log.Printf("Migration: Adding %s column to groups table...", columnName)
			_, err = DB.Exec(fmt.Sprintf("ALTER TABLE groups ADD COLUMN %s %s DEFAULT %s", columnName, columnType, defaultValue))
			if err != nil {
				log.Printf("Failed to migrate groups table (%s): %v", columnName, err)
			} else {
				log.Printf("✅ Migration successful: %s column added", columnName)
			}
		}
	}

	migrateGroupTable("group_rule_mode", "VARCHAR", "'free'")
	migrateGroupTable("thinking_enabled", "BOOLEAN", "FALSE")
	migrateGroupTable("simplified_output", "BOOLEAN", "FALSE")
	migrateGroupTable("custom_rule", "TEXT", "''")

	// Ensure default group exists and contains all agents
	if err := ensureDefaultGroup(); err != nil {
		log.Printf("Failed to ensure default group: %v", err)
	}

	fmt.Println("✅ DuckDB Initialized (Relational + Graph)")
}

func ensureDefaultGroup() error {
	// Check if default group exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM groups WHERE id = 'default'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check default group: %w", err)
	}

	if count == 0 {
		// Create default group
		_, err := DB.Exec("INSERT INTO groups (id, name) VALUES ('default', '所有人 (Default Group)')")
		if err != nil {
			return fmt.Errorf("failed to create default group: %w", err)
		}
		log.Println("✅ Created default group")
	}

	// Get all agent IDs
	rows, err := DB.Query("SELECT id FROM agents")
	if err != nil {
		return fmt.Errorf("failed to fetch agents: %w", err)
	}
	defer rows.Close()

	var agentIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan agent ID: %w", err)
		}
		agentIDs = append(agentIDs, id)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating agents: %w", err)
	}

	// Clear existing members of default group
	_, err = DB.Exec("DELETE FROM group_members WHERE group_id = 'default'")
	if err != nil {
		return fmt.Errorf("failed to clear default group members: %w", err)
	}

	// Add all agents as members
	for _, agentID := range agentIDs {
		_, err := DB.Exec("INSERT INTO group_members (group_id, user_id) VALUES ('default', ?)", agentID)
		if err != nil {
			return fmt.Errorf("failed to add agent %s to default group: %w", agentID, err)
		}
	}

	log.Printf("✅ Default group synced with %d agents", len(agentIDs))
	return nil
}
