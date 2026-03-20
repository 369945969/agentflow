package routes

import (
	"apiServer/db"
	"apiServer/models"
	"net/http"

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
	rows, err := db.DB.Query("SELECT id, name, description, model, system_prompt, created_at FROM agents ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch agents"})
		return
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var a models.Agent
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Model, &a.SystemPrompt, &a.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan agent"})
			return
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

	id := uuid.New().String()
	_, err := db.DB.Exec("INSERT INTO agents (id, name, description, model, system_prompt) VALUES (?, ?, ?, ?, ?)",
		id, input.Name, input.Description, input.Model, input.SystemPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create agent"})
		return
	}

	var newAgent models.Agent
	err = db.DB.QueryRow("SELECT id, name, description, model, system_prompt, created_at FROM agents WHERE id = ?", id).
		Scan(&newAgent.ID, &newAgent.Name, &newAgent.Description, &newAgent.Model, &newAgent.SystemPrompt, &newAgent.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch new agent"})
		return
	}

	c.JSON(http.StatusCreated, newAgent)
}

func deleteAgent(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM agents WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete agent"})
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

	_, err := db.DB.Exec("UPDATE agents SET name = ?, description = ?, model = ?, system_prompt = ? WHERE id = ?",
		input.Name, input.Description, input.Model, input.SystemPrompt, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update agent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent updated"})
}
