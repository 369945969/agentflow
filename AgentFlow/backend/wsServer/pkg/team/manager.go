package team

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// API Server 地址
const APIServerURL = "http://localhost:3000"

// InitTeam 初始化或刷新群组上下文
func (sm *SessionManager) InitTeam(groupID string) (*GroupContext, error) {
	// 1. 获取群组详情 (包含设置和成员)
	// 由于我们的 getGroups 接口返回的是列表，这里我们遍历查找 (或者后端应该提供 getGroupById)
	// 暂时假设我们从 getGroups 过滤
	ctx := sm.GetGroupContext(groupID)
	
	// 从 API 拉取最新数据
	groupData, err := fetchGroupFromAPI(groupID)
	if err != nil {
		return nil, err
	}

	ctx.Mutex.Lock()
	defer ctx.Mutex.Unlock()

	// 2. 更新规则
	ctx.Rules = GroupRules{
		Mode:             groupData.GroupRuleMode,
		ThinkingEnabled:  groupData.ThinkingEnabled,
		SimplifiedOutput: groupData.SimplifiedOutput,
		CustomRule:       groupData.CustomRule,
	}

	// 3. 更新成员画像 & 初始化 Session
	// Main Session (Leader)
	if ctx.MainSessionID == "" {
		// Main Agent 使用默认模型或特定高智商模型
		sess, err := CreateOpenCodeSession("") // 使用默认模型
		if err == nil {
			ctx.MainSessionID = sess
			log.Printf("[Team] Group %s Main Session initialized: %s", groupID, sess)
		}
	}

	// Member Sessions
	newProfiles := make(map[string]AgentProfile)
	for _, m := range groupData.Members {
		// 构造画像
		profile := AgentProfile{
			ID:           m.ID,
			Name:         m.Name,
			Description:  m.Description,
			ModelID:      m.ModelID,
			SystemPrompt: m.SystemPrompt,
			Skills:       m.Skills,
		}
		newProfiles[m.ID] = profile

		// 检查 Session 是否存在
		if _, ok := ctx.MemberSessions[m.ID]; !ok {
			// 为该成员创建 Session，绑定其特定的 Model
			sess, err := CreateOpenCodeSession(m.ModelID)
			if err == nil {
				ctx.MemberSessions[m.ID] = sess
				log.Printf("[Team] Member %s (%s) Session initialized: %s", m.Name, m.ID, sess)
			}
		}
	}
	ctx.Profiles = newProfiles

	return ctx, nil
}

// 辅助结构体用于解析 API 响应
type APIGroupResponse struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	GroupRuleMode    string              `json:"group_rule_mode"`
	ThinkingEnabled  bool                `json:"thinking_enabled"`
	SimplifiedOutput bool                `json:"simplified_output"`
	CustomRule       string              `json:"custom_rule"`
	Members          []APIMemberResponse `json:"members"`
}

type APIMemberResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ModelID      string   `json:"model"` // 注意：API 返回的是 model 字段
	SystemPrompt string   `json:"system_prompt"`
	Skills       []string `json:"skills"` // 假设后端 API 返回了 skills
}

func fetchGroupFromAPI(groupID string) (*APIGroupResponse, error) {
	// 由于后端目前没有 /api/groups/:id 详情接口（只有列表），我们先用列表接口过滤
	// 优化建议：后端应该增加 GET /api/groups/:id
	resp, err := http.Get(fmt.Sprintf("%s/api/groups/", APIServerURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var groups []APIGroupResponse
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}

	for _, g := range groups {
		if g.ID == groupID {
			// 还需要获取 Member 的详细信息 (Description, SystemPrompt)
			// 因为 groups 接口返回的 member 可能信息不全
			// 这里做一个简单的模拟：遍历成员 ID 再去调 /api/agents/:id 获取详情
			// 为了性能，实际生产中应该在 Group 详情接口一次性返回
			for i, m := range g.Members {
				agentInfo, err := fetchAgentDetail(m.ID)
				if err == nil {
					g.Members[i] = *agentInfo
				}
			}
			return &g, nil
		}
	}
	return nil, fmt.Errorf("group not found: %s", groupID)
}

func fetchAgentDetail(agentID string) (*APIMemberResponse, error) {
    // 假设后端支持 GET /api/agents/:id，但目前的 getAgents 是返回列表
    // 我们暂时从列表获取
    resp, err := http.Get(fmt.Sprintf("%s/api/agents/", APIServerURL))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var agents []APIMemberResponse
    body, _ := io.ReadAll(resp.Body)
    json.Unmarshal(body, &agents)
    
    for _, a := range agents {
        if a.ID == agentID {
            return &a, nil
        }
    }
    return nil, fmt.Errorf("agent not found")
}

// RouterDecision 路由决策结果
type RouterDecision struct {
	SelectedAgentID string `json:"selected_agent_id"`
	Reason          string `json:"reason"`
	Instruction     string `json:"instruction"`
}

// DecideNextAgent 调用 Main Agent 进行决策
// 注意：这是一个同步阻塞调用，实际应该处理流式返回并解析 JSON
func (ctx *GroupContext) DecideNextAgent(userMessage string) (*RouterDecision, error) {
	// 1. 构建 Prompt
	profiles := make([]AgentProfile, 0, len(ctx.Profiles))
	for _, p := range ctx.Profiles {
		profiles = append(profiles, p)
	}
	systemPrompt := BuildRouterPrompt(ctx.Rules, profiles)
	
	// 2. 发送请求给 Main Session
	// 这里我们需要一种机制来发送 System Prompt + User Message，并获取完整回复
	// 由于 wsServer 目前只转发流，这里我们需要一个非流式的 Helper 或者自行收集流
	// 简化起见，我们构造一个特殊的 Prompt
	
	fullPrompt := fmt.Sprintf("%s\n\n【用户消息】：%s", systemPrompt, userMessage)
	
	req, err := BuildRequest(ctx.MainSessionID, "", fullPrompt) // Main Agent 用默认模型
	if err != nil {
		return nil, err
	}
	
	// 发送并等待结果
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 解析 SSE 流，拼接完整内容
	// ... (这里需要复用流解析逻辑，但只为了获取最终 JSON)
	// 为简化演示，假设 OpenCode 支持非流式返回，或者我们在这里简单读取
	// 实际代码需要完整的 SSE Parser
	
	return &RouterDecision{
		SelectedAgentID: "mock_agent_id", // 这里的解析逻辑比较复杂，下一步在 main.go 实现
		Reason: "Mock reason",
	}, nil
}
