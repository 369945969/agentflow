package team

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// AgentProfile 定义群成员的画像
type AgentProfile struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ModelID      string   `json:"model_id"`
	SystemPrompt string   `json:"system_prompt"`
	Skills       []string `json:"skills"`
}

// GroupContext 保存一个群组的运行时状态
type GroupContext struct {
	GroupID        string
	MainSessionID  string                    // 主控 Agent 的 Session ID
	MemberSessions map[string]string         // UserID -> SessionID
	Profiles       map[string]AgentProfile   // UserID -> Profile
	Rules          GroupRules                // 群组规则配置
	Mutex          sync.RWMutex
}

type GroupRules struct {
	Mode             string // free, expert, custom
	ThinkingEnabled  bool
	SimplifiedOutput bool
	CustomRule       string
}

// SessionManager 管理所有群组的上下文
type SessionManager struct {
	groups map[string]*GroupContext
	mutex  sync.RWMutex
}

var Manager = &SessionManager{
	groups: make(map[string]*GroupContext),
}

func (sm *SessionManager) GetGroupContext(groupID string) *GroupContext {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.groups[groupID]; !exists {
		sm.groups[groupID] = &GroupContext{
			GroupID:        groupID,
			MemberSessions: make(map[string]string),
			Profiles:       make(map[string]AgentProfile),
		}
	}
	return sm.groups[groupID]
}

// OpenCodeClient 封装与 OpenCode 的 HTTP 交互
const OpenCodeBaseURL = "http://localhost:3001"

func CreateOpenCodeSession(modelID string) (string, error) {
	// OpenCode 的 session 创建通常是隐式的，或者通过首次请求建立
	// 这里我们生成一个唯一的 Session ID 并尝试通过 /api/session/init (假设有) 或直接使用
	// 为了适配现有 OpenCode，我们直接生成 UUID 作为 SessionID
	// 实际应用中，这里可能需要调用 OpenCode 的 POST /api/session
	return fmt.Sprintf("sess_%d", timeNow()), nil
}

// 简单的辅助函数获取当前纳秒
func timeNow() int64 {
	return time.Now().UnixNano()
}

// SendToOpenCode 发送消息并获取流式响应 (这里只做简单的请求封装，流式处理在 main.go 中)
func BuildRequest(sessionID string, modelID string, message string) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/sse/chat", OpenCodeBaseURL)
	payload := map[string]string{
		"message": message,
	}
	body, _ := json.Marshal(payload)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", sessionID)
	
	// 如果指定了 ModelID，通过 Header 传递给 DuckDB Router 插件
	if modelID != "" {
		req.Header.Set("X-Model-ID", modelID)
	}
	
	return req, nil
}
