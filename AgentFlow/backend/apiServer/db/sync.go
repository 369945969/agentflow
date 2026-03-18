package db

import (
	"apiServer/config"
	"bufio"
	"fmt"
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

	// 1. Clear skills table
	_, err := DB.Exec("DELETE FROM skills")
	if err != nil {
		log.Fatalf("Failed to clear skills table: %v", err)
	}
	log.Println("🧹 Cleared skills table")

	// 2. Scan skills directory
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

		// 3. Parse SKILL.md for name and description
		name, description := parseSkillMd(skillMdPath)
		if name == "" {
			name = entry.Name() // Fallback to directory name
		}

		skillList = append(skillList, SkillInfo{
			Name:        name,
			Description: description,
			DirName:     entry.Name(),
		})
	}

	// 4. Sort skills by Name alphabetically
	sort.Slice(skillList, func(i, j int) bool {
		return strings.ToLower(skillList[i].Name) < strings.ToLower(skillList[j].Name)
	})

	// 5. Insert into database in sorted order
	for _, s := range skillList {
		id := uuid.New().String()
		absPath, _ := filepath.Abs(filepath.Join(skillsPath, s.DirName))
		_, err = DB.Exec("INSERT INTO skills (id, name, description, type, version, icon) VALUES (?, ?, ?, ?, ?, ?)",
			id, s.Name, s.Description, "Local", absPath, "lucide:terminal")
		if err != nil {
			log.Printf("Failed to insert skill %s: %v", s.Name, err)
			continue
		}
		log.Printf("✅ Synced skill: %s", s.Name)
	}

	fmt.Println("🚀 Skills sync completed")
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
