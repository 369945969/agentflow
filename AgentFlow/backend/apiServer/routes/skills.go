package routes

import (
	"apiServer/db"
	"apiServer/models"
	"database/sql"
	"log"
	"net/http"

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

	var skills []models.Skill
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
	var input struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	skillData := models.Skill{
		ID:          id,
		Name:        "Remote Skill",
		Type:        "HTTP API",
		Description: "Installed from " + input.URL,
		Version:     "v1.0.0",
		Icon:        "lucide:link",
		Logic:       "# Fetched from " + input.URL,
	}

	_, err := db.DB.Exec("INSERT INTO skills (id, name, type, description, version, icon, logic) VALUES (?, ?, ?, ?, ?, ?, ?)",
		skillData.ID, skillData.Name, skillData.Type, skillData.Description, skillData.Version, skillData.Icon, skillData.Logic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Installation failed"})
		return
	}

	var newSkill models.Skill
	var color, logic sql.NullString
	err = db.DB.QueryRow("SELECT id, name, type, description, version, icon, color, logic, created_at FROM skills WHERE id = ?", id).
		Scan(&newSkill.ID, &newSkill.Name, &newSkill.Type, &newSkill.Description, &newSkill.Version, &newSkill.Icon, &color, &logic, &newSkill.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch installed skill"})
		return
	}
	newSkill.Color = color.String
	newSkill.Logic = logic.String

	c.JSON(http.StatusOK, newSkill)
}

func deleteSkill(c *gin.Context) {
	id := c.Param("id")

	// Delete from edges first (mirroring TypeScript logic)
	// Note: The TS logic used source_id and target_id which weren't in the schema from duckdb.ts, but let's follow the schema from TS code
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
