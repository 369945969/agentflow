package routes

import (
	"apiServer/db"
	"apiServer/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterModelRoutes(r *gin.Engine) {
	group := r.Group("/api/models")
	{
		group.GET("/", getModels)
		group.POST("/", createModel)
		group.PUT("/:id", updateModel)
		group.DELETE("/:id", deleteModel)
	}
}

func getModels(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, provider, base_url, api_key, model_name, description, is_default, created_at FROM models ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Failed to fetch models: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch models: " + err.Error()})
		return
	}
	defer rows.Close()

	modelsList := []models.Model{}
	for rows.Next() {
		var m models.Model
		if err := rows.Scan(&m.ID, &m.Name, &m.Provider, &m.BaseURL, &m.APIKey, &m.ModelName, &m.Description, &m.IsDefault, &m.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan model"})
			return
		}
		modelsList = append(modelsList, m)
	}
	c.JSON(http.StatusOK, modelsList)
}

func createModel(c *gin.Context) {
	var input models.Model
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.IsDefault {
		_, err := db.DB.Exec("UPDATE models SET is_default = FALSE")
		if err != nil {
			log.Printf("Failed to reset defaults: %v", err)
		}
	}

	input.ID = uuid.New().String()
	_, err := db.DB.Exec("INSERT INTO models (id, name, provider, base_url, api_key, model_name, description, is_default) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		input.ID, input.Name, input.Provider, input.BaseURL, input.APIKey, input.ModelName, input.Description, input.IsDefault)
	if err != nil {
		log.Printf("Failed to create model: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func updateModel(c *gin.Context) {
	id := c.Param("id")
	var input models.Model
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.IsDefault {
		_, err := db.DB.Exec("UPDATE models SET is_default = FALSE")
		if err != nil {
			log.Printf("Failed to reset defaults: %v", err)
		}
	}

	_, err := db.DB.Exec("UPDATE models SET name=?, provider=?, base_url=?, api_key=?, model_name=?, description=?, is_default=? WHERE id=?",
		input.Name, input.Provider, input.BaseURL, input.APIKey, input.ModelName, input.Description, input.IsDefault, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Model updated"})
}

func deleteModel(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM models WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete model"})
		return
	}
	c.Status(http.StatusNoContent)
}
