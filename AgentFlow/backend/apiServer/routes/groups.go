package routes

import (
	"apiServer/db"
	"apiServer/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterGroupRoutes(r *gin.Engine) {
	group := r.Group("/api/groups")
	{
		group.POST("/", createGroup)
		group.GET("/", getGroups)
		group.PUT("/:id", updateGroupSettings)
		group.DELETE("/:id", deleteGroup)
	}
}

type CreateGroupInput struct {
	Name    string   `json:"name" binding:"required"`
	Members []string `json:"members"`
}

type UpdateGroupSettingsInput struct {
	GroupRuleMode    string `json:"group_rule_mode"`
	ThinkingEnabled  bool   `json:"thinking_enabled"`
	SimplifiedOutput bool   `json:"simplified_output"`
	CustomRule       string `json:"custom_rule"`
}

func createGroup(c *gin.Context) {
	var input CreateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	newGroup := models.Group{
		ID:   uuid.New().String(),
		Name: input.Name,
	}
	_, err = tx.Exec("INSERT INTO groups (id, name) VALUES (?, ?)", newGroup.ID, newGroup.Name)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	for _, memberID := range input.Members {
		_, err = tx.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", newGroup.ID, memberID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add members to group"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	newGroup.Members = input.Members
	c.JSON(http.StatusCreated, newGroup)
}

func getGroups(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT 
			g.id, 
			g.name, 
			g.created_at, 
			g.group_rule_mode,
			g.thinking_enabled,
			g.simplified_output,
			g.custom_rule,
			gm.user_id,
			a.name as agent_name
		FROM groups g 
		LEFT JOIN group_members gm ON g.id = gm.group_id 
		LEFT JOIN agents a ON gm.user_id = a.id 
		ORDER BY g.id = 'default' DESC, g.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}
	defer rows.Close()

	groupsMap := make(map[string]map[string]interface{})
	for rows.Next() {
		var groupID, groupName, createdAt, groupRuleMode, customRule string
		var thinkingEnabled, simplifiedOutput bool
		var memberID, agentName sql.NullString
		if err := rows.Scan(&groupID, &groupName, &createdAt, &groupRuleMode, &thinkingEnabled, &simplifiedOutput, &customRule, &memberID, &agentName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan group"})
			return
		}

		group, exists := groupsMap[groupID]
		if !exists {
			group = map[string]interface{}{
				"id":                groupID,
				"name":              groupName,
				"lastMsg":           "欢迎加入群聊",
				"created_at":        createdAt,
				"group_rule_mode":   groupRuleMode,
				"thinking_enabled":  thinkingEnabled,
				"simplified_output": simplifiedOutput,
				"custom_rule":       customRule,
				"members":           []map[string]string{},
			}
		}

		if memberID.Valid {
			member := map[string]string{
				"id":   memberID.String,
				"name": agentName.String,
			}
			// 如果agent_name为空，使用memberID作为后备
			if agentName.String == "" {
				member["name"] = memberID.String
			}
			group["members"] = append(group["members"].([]map[string]string), member)
		}
		groupsMap[groupID] = group
	}

	// 处理查询错误
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while iterating rows"})
		return
	}

	// 转换为列表
	groupsList := make([]map[string]interface{}, 0)
	for _, group := range groupsMap {
		groupsList = append(groupsList, group)
	}

	c.JSON(http.StatusOK, groupsList)
}

func updateGroupSettings(c *gin.Context) {
	id := c.Param("id")

	// Prevent updating default group settings? Allow updates.

	var input UpdateGroupSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate group_rule_mode
	validModes := map[string]bool{"free": true, "expert": true, "custom": true}
	if !validModes[input.GroupRuleMode] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group_rule_mode, must be 'free', 'expert', or 'custom'"})
		return
	}

	// Update the group settings
	_, err := db.DB.Exec(`
		UPDATE groups 
		SET group_rule_mode = ?, 
		    thinking_enabled = ?, 
		    simplified_output = ?, 
		    custom_rule = ?
		WHERE id = ?`,
		input.GroupRuleMode,
		input.ThinkingEnabled,
		input.SimplifiedOutput,
		input.CustomRule,
		id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group settings updated"})
}

func deleteGroup(c *gin.Context) {
	id := c.Param("id")

	// Prevent deletion of default group
	if id == "default" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete default group"})
		return
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Delete group members first
	_, err = tx.Exec("DELETE FROM group_members WHERE group_id = ?", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group members"})
		return
	}

	// Delete the group
	_, err = tx.Exec("DELETE FROM groups WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}
