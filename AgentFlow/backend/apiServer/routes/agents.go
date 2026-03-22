package routes

import (
	"apiServer/db"
	"apiServer/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterAgentRoutes(r *gin.Engine) {
	agentGroup := r.Group("/api/agents")
	{
		agentGroup.GET("/", getAgents)
		agentGroup.POST("/", createAgent)
		agentGroup.PUT("/:id", updateAgent)
		agentGroup.DELETE("/:id", deleteAgent)
	}
}

func getAgents(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, description, model, system_prompt, thinking_enabled, simplified_output, created_at FROM agents ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch agents"})
		return
	}
	defer rows.Close()

	agents := []models.Agent{}
	for rows.Next() {
		var a models.Agent
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Model, &a.SystemPrompt, &a.ThinkingEnabled, &a.SimplifiedOutput, &a.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan agent"})
			return
		}

		// Fetch skills for this agent
		skillRows, err := db.DB.Query("SELECT skill_id FROM agent_skills WHERE agent_id = ?", a.ID)
		if err == nil {
			var skills []string
			for skillRows.Next() {
				var skillID string
				if err := skillRows.Scan(&skillID); err == nil {
					skills = append(skills, skillID)
				}
			}
			skillRows.Close()
			a.Skills = skills
		}

		agents = append(agents, a)
	}

	c.JSON(http.StatusOK, agents)
}

func createAgent(c *gin.Context) {
	var input models.Agent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	id := uuid.New().String()
	_, err = tx.Exec("INSERT INTO agents (id, name, description, model, system_prompt, thinking_enabled, simplified_output) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, input.Name, input.Description, input.Model, input.SystemPrompt, input.ThinkingEnabled, input.SimplifiedOutput)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create agent"})
		return
	}

	// Add agent to default group
	_, err = tx.Exec("INSERT INTO group_members (group_id, user_id) VALUES ('default', ?)", id)
	if err != nil {
		log.Printf("Failed to add agent to default group: %v", err)
	}

	// Add skills
	for _, skillID := range input.Skills {
		_, err = tx.Exec("INSERT INTO agent_skills (agent_id, skill_id) VALUES (?, ?)", id, skillID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add skills to agent"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	var newAgent models.Agent
	err = db.DB.QueryRow("SELECT id, name, description, model, system_prompt, thinking_enabled, simplified_output, created_at FROM agents WHERE id = ?", id).
		Scan(&newAgent.ID, &newAgent.Name, &newAgent.Description, &newAgent.Model, &newAgent.SystemPrompt, &newAgent.ThinkingEnabled, &newAgent.SimplifiedOutput, &newAgent.CreatedAt)
	if err == nil {
		newAgent.Skills = input.Skills
	}

	c.JSON(http.StatusCreated, newAgent)
}

func deleteAgent(c *gin.Context) {
	id := c.Param("id")
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	_, err = tx.Exec("DELETE FROM agents WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete agent"})
		return
	}

	tx.Exec("DELETE FROM group_members WHERE user_id = ?", id)
	tx.Exec("DELETE FROM agent_skills WHERE agent_id = ?", id)

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

func updateAgent(c *gin.Context) {
	id := c.Param("id")
	var input models.Agent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	_, err = tx.Exec("UPDATE agents SET name = ?, description = ?, model = ?, system_prompt = ?, thinking_enabled = ?, simplified_output = ? WHERE id = ?",
		input.Name, input.Description, input.Model, input.SystemPrompt, input.ThinkingEnabled, input.SimplifiedOutput, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update agent"})
		return
	}

	// Update skills: delete old and insert new
	if _, err := tx.Exec("DELETE FROM agent_skills WHERE agent_id = ?", id); err != nil {
		tx.Rollback()
		log.Printf("Failed to delete existing agent skills for agent %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing agent skills"})
		return
	}

	seenSkills := make(map[string]struct{}, len(input.Skills))
	for _, skillID := range input.Skills {
		skillID = strings.TrimSpace(skillID)
		if skillID == "" {
			continue
		}
		if _, exists := seenSkills[skillID]; exists {
			continue
		}
		seenSkills[skillID] = struct{}{}

		_, err = tx.Exec("INSERT OR IGNORE INTO agent_skills (agent_id, skill_id) VALUES (?, ?)", id, skillID)
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to insert agent skill for agent %s skill %s: %v", id, skillID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update agent skills"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent updated"})
}
