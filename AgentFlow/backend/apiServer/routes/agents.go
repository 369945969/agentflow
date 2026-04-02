package routes

import (
	"apiServer/db"
	"apiServer/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		agentGroup.POST("/ai-generate", aiGenerateAgent)
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

type AIGenerateInput struct {
	Prompt string `json:"prompt" binding:"required"`
}

func aiGenerateAgent(c *gin.Context) {
	var input AIGenerateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Fetch default model
	var model models.Model
	err := db.DB.QueryRow("SELECT id, provider, base_url, api_key, model_name FROM models WHERE is_default = TRUE LIMIT 1").
		Scan(&model.ID, &model.Provider, &model.BaseURL, &model.APIKey, &model.ModelName)
	
	if err != nil {
		err = db.DB.QueryRow("SELECT id, provider, base_url, api_key, model_name FROM models LIMIT 1").
			Scan(&model.ID, &model.Provider, &model.BaseURL, &model.APIKey, &model.ModelName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No model configured in system"})
			return
		}
	}

	// 2. Fetch available skills
	rows, _ := db.DB.Query("SELECT id, name, description FROM skills")
	var availableSkills []map[string]string
	if rows != nil {
		for rows.Next() {
			var id, name, desc string
			rows.Scan(&id, &name, &desc)
			availableSkills = append(availableSkills, map[string]string{"id": id, "name": name, "description": desc})
		}
		rows.Close()
	}
	skillsJSON, _ := json.Marshal(availableSkills)

	// 3. Prepare Prompt (Fixed backtick issue)
	systemPrompt := fmt.Sprintf(`你是 AI 智能体架构专家。根据用户的需求描述，设计一个最适合的 AI 智能体配置。
可供选择的现有技能列表: %s

你必须直接输出且仅输出一个合法的 JSON 对象，不要包含任何解释、不要包含 Markdown 标记。
JSON 结构如下:
{
  "name": "智能体显示名称",
  "description": "一句话职责描述",
  "system_prompt": "详细的系统提示词，包含角色设定、工作流程、回复规范等",
  "skills": ["从技能列表中选择的技能ID数组"]
}`, string(skillsJSON))

	// 4. Call LLM
	apiURL := fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(model.BaseURL, "/"))
	if (model.Provider == "Ollama" || strings.Contains(strings.ToLower(model.BaseURL), "localhost")) && !strings.Contains(model.BaseURL, "/v1") {
		apiURL = fmt.Sprintf("%s/v1/chat/completions", strings.TrimSuffix(model.BaseURL, "/"))
	}

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": model.ModelName,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": input.Prompt},
		},
		"temperature": 0.3,
	})

	log.Printf("[AI-Generate] Calling LLM at %s with model %s", apiURL, model.ModelName)

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	if model.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+model.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call LLM: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("[AI-Generate] LLM returned status %d: %s", resp.StatusCode, string(respBody))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("LLM returned error %d", resp.StatusCode)})
		return
	}

	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse LLM response"})
		return
	}

	if len(openAIResp.Choices) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LLM returned empty choices"})
		return
	}

	content := strings.TrimSpace(openAIResp.Choices[0].Message.Content)
	log.Printf("[AI-Generate] LLM raw output: %s", content)

	// Clean JSON if LLM added markdown fences
	if strings.HasPrefix(content, "```") {
		lines := strings.Split(content, "\n")
		var cleanLines []string
		for _, line := range lines {
			if !strings.HasPrefix(line, "```") {
				cleanLines = append(cleanLines, line)
			}
		}
		content = strings.Join(cleanLines, "\n")
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// Try to find JSON block manually
		start := strings.Index(content, "{")
		end := strings.LastIndex(content, "}")
		if start >= 0 && end > start {
			jsonStr := content[start : end+1]
			if err := json.Unmarshal([]byte(jsonStr), &result); err == nil {
				c.JSON(http.StatusOK, result)
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LLM output is not valid JSON", "raw": content})
		return
	}

	c.JSON(http.StatusOK, result)
}
