package routes

import (
	"apiServer/config"
	"apiServer/db"
	"apiServer/models"
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterSkillRoutes(r *gin.Engine) {
	skillGroup := r.Group("/api/skills")
	{
		skillGroup.GET("/", getSkills)
		skillGroup.POST("/", createSkill)
		skillGroup.POST("/install", installSkill)
		skillGroup.DELETE("/:id", deleteSkill)
	}
}

func getSkills(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, type, description, version, icon, color, logic, created_at FROM skills ORDER BY created_at ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}
	defer rows.Close()

	skills := []models.Skill{}
	for rows.Next() {
		var s models.Skill
		var color, logic sql.NullString
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Description, &s.Version, &s.Icon, &color, &logic, &s.CreatedAt); err != nil {
			log.Printf("Error scanning skill: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan skill"})
			return
		}
		s.Color = color.String
		s.Logic = logic.String
		skills = append(skills, s)
	}

	c.JSON(http.StatusOK, skills)
}

func createSkill(c *gin.Context) {
	var input models.Skill
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	_, err := db.DB.Exec("INSERT INTO skills (id, name, type, description, logic, version, icon, color) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		id, input.Name, input.Type, input.Description, input.Logic, input.Version, input.Icon, input.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create skill"})
		return
	}

	var newSkill models.Skill
	var color, logic sql.NullString
	err = db.DB.QueryRow("SELECT id, name, type, description, version, icon, color, logic, created_at FROM skills WHERE id = ?", id).
		Scan(&newSkill.ID, &newSkill.Name, &newSkill.Type, &newSkill.Description, &newSkill.Version, &newSkill.Icon, &color, &logic, &newSkill.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch new skill"})
		return
	}
	newSkill.Color = color.String
	newSkill.Logic = logic.String

	c.JSON(http.StatusCreated, newSkill)
}

func installSkill(c *gin.Context) {
	// 1. Get configuration
	skillsPath := config.AppConfig.Skills.Path
	if skillsPath == "" {
		skillsPath = "../../skills" // Default fallback
	}

	var tempZipPath string
	var cleanup func()

	// 2. Handle either File upload or URL
	file, err := c.FormFile("file")
	if err == nil {
		// Local file upload
		id := uuid.New().String()
		tempZipPath = filepath.Join(os.TempDir(), id+".zip")
		if err := c.SaveUploadedFile(file, tempZipPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
			return
		}
		cleanup = func() { os.Remove(tempZipPath) }
	} else {
		// URL install
		url := c.PostForm("url")
		if url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Neither file nor URL provided"})
			return
		}

		id := uuid.New().String()
		tempZipPath = filepath.Join(os.TempDir(), id+".zip")
		if err := downloadFile(url, tempZipPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download from URL: " + err.Error()})
			return
		}
		cleanup = func() { os.Remove(tempZipPath) }
	}
	defer cleanup()

	// 3. Unzip and Process
	extractedDir, err := unzipSkill(tempZipPath, skillsPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unzip skill: " + err.Error()})
		return
	}

	// 4. Parse metadata (_meta.json or package.json)
	skillInfo, err := parseSkillMetadata(extractedDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid skill package: " + err.Error()})
		return
	}

	// 5. Register in database
	id := filepath.Base(extractedDir)
	_, err = db.DB.Exec("INSERT OR IGNORE INTO skills (id, name, type, description, version, icon, color) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, skillInfo.Name, skillInfo.Type, skillInfo.Description, skillInfo.Version, skillInfo.Icon, skillInfo.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register skill in database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill installed successfully", "id": id})
}

func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzipSkill(src string, destBase string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// Determine the root directory name inside the zip
	var rootDir string
	if len(r.File) > 0 {
		parts := strings.Split(r.File[0].Name, "/")
		rootDir = parts[0]
	}

	if rootDir == "" {
		return "", fmt.Errorf("empty zip file")
	}

	dest := filepath.Join(destBase, rootDir)

	for _, f := range r.File {
		fpath := filepath.Join(destBase, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return "", err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}
	}

	return dest, nil
}

type SkillMeta struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

func parseSkillMetadata(dir string) (SkillMeta, error) {
	var meta SkillMeta

	// Try _meta.json first
	metaPath := filepath.Join(dir, "_meta.json")
	if _, err := os.Stat(metaPath); err == nil {
		data, _ := os.ReadFile(metaPath)
		json.Unmarshal(data, &meta)
	} else {
		// Try package.json
		pkgPath := filepath.Join(dir, "package.json")
		if _, err := os.Stat(pkgPath); err == nil {
			data, _ := os.ReadFile(pkgPath)
			json.Unmarshal(data, &meta)
		} else {
			// Minimal fallback
			meta.Name = filepath.Base(dir)
			meta.Type = "Plugin"
		}
	}

	if meta.Name == "" {
		meta.Name = filepath.Base(dir)
	}
	return meta, nil
}

func deleteSkill(c *gin.Context) {
	id := c.Param("id")

	_, err := db.DB.Exec("DELETE FROM edges WHERE from_id = ? OR to_id = ?", id, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related edges"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM skills WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill"})
		return
	}

	c.Status(http.StatusNoContent)
}
