package routes

import (
	"apiServer/db"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Workflow struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Graph       json.RawMessage `json:"graph"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func RegisterWorkflowRoutes(r *gin.Engine) {
	group := r.Group("/api/workflows")
	{
		group.GET("/", getWorkflows)
		group.GET("/:id", getWorkflow)
		group.POST("/", createWorkflow)
		group.PUT("/:id", updateWorkflow)
		group.DELETE("/:id", deleteWorkflow)
	}
}

func getWorkflows(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, description, graph, created_at, updated_at FROM workflows ORDER BY updated_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workflows"})
		return
	}
	defer rows.Close()

	workflows := []Workflow{}
	for rows.Next() {
		var w Workflow
		var graphData interface{}
		if err := rows.Scan(&w.ID, &w.Name, &w.Description, &graphData, &w.CreatedAt, &w.UpdatedAt); err != nil {
			log.Printf("Scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan workflow"})
			return
		}
		
		// Convert the interfaced map back to json.RawMessage
		jsonBytes, _ := json.Marshal(graphData)
		w.Graph = json.RawMessage(jsonBytes)
		
		workflows = append(workflows, w)
	}
	c.JSON(http.StatusOK, workflows)
}

func getWorkflow(c *gin.Context) {
	id := c.Param("id")
	var w Workflow
	var graphData interface{}
	err := db.DB.QueryRow("SELECT id, name, description, graph, created_at, updated_at FROM workflows WHERE id = ?", id).
		Scan(&w.ID, &w.Name, &w.Description, &graphData, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}
	
	jsonBytes, _ := json.Marshal(graphData)
	w.Graph = json.RawMessage(jsonBytes)
	
	c.JSON(http.StatusOK, w)
}

func createWorkflow(c *gin.Context) {
	var input Workflow
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = uuid.New().String()
	now := time.Now()
	_, err := db.DB.Exec("INSERT INTO workflows (id, name, description, graph, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		input.ID, input.Name, input.Description, string(input.Graph), now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func updateWorkflow(c *gin.Context) {
	id := c.Param("id")
	var input Workflow
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec("UPDATE workflows SET name = ?, description = ?, graph = ?, updated_at = ? WHERE id = ?",
		input.Name, input.Description, string(input.Graph), time.Now(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workflow updated"})
}

func deleteWorkflow(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM workflows WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workflow"})
		return
	}
	c.Status(http.StatusNoContent)
}
