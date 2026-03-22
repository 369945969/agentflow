package db

import (
	"apiServer/config"
	"bufio"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type SkillInfo struct {
	Name        string
	Description string
	DirName     string
	AbsPath     string
}

type SkillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func SyncSkills() {
	skillsPath := config.AppConfig.Skills.Path
	if skillsPath == "" {
		log.Println("⚠️ Skills path not configured, skipping sync")
		return
	}

	// 1. Scan skills directory
	entries, err := os.ReadDir(skillsPath)
	if err != nil {
		log.Fatalf("Failed to read skills directory: %v", err)
	}

	var skillList []SkillInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillDir := filepath.Join(skillsPath, entry.Name())
		skillMdPath := filepath.Join(skillDir, "SKILL.md")

		if _, err := os.Stat(skillMdPath); os.IsNotExist(err) {
			continue
		}

		// 2. Parse SKILL.md for name and description
		name, description := parseSkillMd(skillMdPath)
		if name == "" {
			name = entry.Name() // Fallback to directory name
		}
		absPath, err := filepath.Abs(skillDir)
		if err != nil {
			log.Printf("Failed to resolve absolute path for skill %s: %v", entry.Name(), err)
			continue
		}

		skillList = append(skillList, SkillInfo{
			Name:        name,
			Description: description,
			DirName:     entry.Name(),
			AbsPath:     absPath,
		})
	}

	// 3. Sort skills by Name alphabetically
	sort.Slice(skillList, func(i, j int) bool {
		return strings.ToLower(skillList[i].Name) < strings.ToLower(skillList[j].Name)
	})

	tx, err := DB.Begin()
	if err != nil {
		log.Fatalf("Failed to start skill sync transaction: %v", err)
	}
	defer tx.Rollback()

	// 4. Upsert into database in sorted order using path + name as the identity.
	for _, s := range skillList {
		if err := upsertLocalSkill(tx, s); err != nil {
			log.Printf("Failed to sync skill %s (%s): %v", s.Name, s.AbsPath, err)
			continue
		}
		log.Printf("✅ Synced skill: %s", s.Name)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit skill sync transaction: %v", err)
	}

	log.Println("🚀 Skills sync completed")
}

func upsertLocalSkill(tx *sql.Tx, skill SkillInfo) error {
	rows, err := tx.Query(
		"SELECT id FROM skills WHERE name = ? AND version = ? ORDER BY created_at ASC, id ASC",
		skill.Name, skill.AbsPath,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var existingIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		existingIDs = append(existingIDs, id)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(existingIDs) == 0 {
		_, err = tx.Exec(
			"INSERT INTO skills (id, name, description, type, version, icon) VALUES (?, ?, ?, ?, ?, ?)",
			uuid.New().String(), skill.Name, skill.Description, "Local", skill.AbsPath, "lucide:terminal",
		)
		return err
	}

	canonicalID := existingIDs[0]
	_, err = tx.Exec(
		"UPDATE skills SET description = ?, type = ?, version = ?, icon = ? WHERE id = ?",
		skill.Description, "Local", skill.AbsPath, "lucide:terminal", canonicalID,
	)
	if err != nil {
		return err
	}

	for _, duplicateID := range existingIDs[1:] {
		if _, err := tx.Exec(
			"INSERT OR IGNORE INTO agent_skills (agent_id, skill_id) SELECT agent_id, ? FROM agent_skills WHERE skill_id = ?",
			canonicalID, duplicateID,
		); err != nil {
			return err
		}
		if _, err := tx.Exec("DELETE FROM agent_skills WHERE skill_id = ?", duplicateID); err != nil {
			return err
		}
		if _, err := tx.Exec("DELETE FROM skills WHERE id = ?", duplicateID); err != nil {
			return err
		}
		log.Printf("🧹 Merged duplicate skill %s (%s) into %s", skill.Name, skill.AbsPath, canonicalID)
	}

	return nil
}

func parseSkillMd(path string) (string, string) {
	file, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	var frontmatterRaw strings.Builder
	scanner := bufio.NewScanner(file)
	inFrontmatter := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				break // End of frontmatter
			}
		}
		if inFrontmatter {
			frontmatterRaw.WriteString(line + "\n")
		}
	}

	var fm SkillFrontmatter
	err = yaml.Unmarshal([]byte(frontmatterRaw.String()), &fm)
	if err != nil || (fm.Name == "" && fm.Description == "") {
		// If YAML parsing fails or returns empty, try simple search
		return searchSimpleMetadata(path)
	}

	return fm.Name, fm.Description
}

func searchSimpleMetadata(path string) (string, string) {
	file, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	var name, description string
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() && count < 20 {
		line := scanner.Text()
		if strings.HasPrefix(strings.ToLower(line), "name:") {
			name = strings.TrimSpace(line[5:])
		} else if strings.HasPrefix(strings.ToLower(line), "description:") {
			description = strings.TrimSpace(line[12:])
		}
		count++
	}
	return name, description
}
