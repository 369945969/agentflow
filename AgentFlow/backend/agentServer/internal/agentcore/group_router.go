package agentcore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"agentServer/internal/protocol"
)

type GroupSessionStore struct {
	mu             sync.RWMutex
	mainSessions   map[string]string            // GroupID -> Main SessionID
	memberSessions map[string]map[string]string // GroupID -> (UserID -> SessionID)
}

func NewGroupSessionStore() *GroupSessionStore {
	return &GroupSessionStore{
		mainSessions:   make(map[string]string),
		memberSessions: make(map[string]map[string]string),
	}
}

func (s *GroupSessionStore) GetMainSession(groupID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mainSessions[groupID]
}

func (s *GroupSessionStore) SetMainSession(groupID, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mainSessions[groupID] = sessionID
}

func (s *GroupSessionStore) GetMemberSession(groupID, userID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if ms, ok := s.memberSessions[groupID]; ok {
		return ms[userID]
	}
	return ""
}

func (s *GroupSessionStore) SetMemberSession(groupID, userID, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.memberSessions[groupID] == nil {
		s.memberSessions[groupID] = make(map[string]string)
	}
	s.memberSessions[groupID][userID] = sessionID
}

type GroupRouter struct {
	core  *AgentCore
	store *GroupSessionStore
}

func NewGroupRouter(core *AgentCore) *GroupRouter {
	return &GroupRouter{
		core:  core,
		store: NewGroupSessionStore(),
	}
}

func (gr *GroupRouter) Handle(ctx context.Context, execCtx ExecutionContext) {
	groupInfo := execCtx.Session.GroupInfo
	if groupInfo == nil || groupInfo.GroupID == "" {
		fallback := execCtx
		fallback.SystemPrompt = buildSingleSystemPrompt(execCtx.UserProfile.SystemPrompt, execCtx.EffectiveConfig.EnableThinking)
		if a, ok := gr.core.pool.Default(); ok {
			fallback.TargetAgent = a
		} else {
			gr.sendError(execCtx, "no_available_agent", "No agent configured")
			return
		}
		go gr.core.Stream(ctx, fallback)
		return
	}

	groupID := groupInfo.GroupID

	// 1) MainAgent Session
	mainSessionID := gr.store.GetMainSession(groupID)
	mainSession, err := gr.core.sessionMgr.GetOrCreateSession(groupInfo.CreatorID, mainSessionID, "group")
	if err != nil {
		gr.sendError(execCtx, "session_error", "Failed to get main session")
		return
	}
	gr.store.SetMainSession(groupID, mainSession.SessionID)

	// 2) Collect allowed members
	var availableMembers []GroupMember
	for _, m := range groupInfo.Members {
		if !m.IsActive {
			continue
		}
		availableMembers = append(availableMembers, m)
	}

	// 3) Orchestrate the conversation loop natively
	if len(availableMembers) == 0 {
		fallback := execCtx
		fallback.SystemPrompt = buildSingleSystemPrompt(execCtx.UserProfile.SystemPrompt, execCtx.EffectiveConfig.EnableThinking)
		if a, ok := gr.core.pool.Default(); ok {
			fallback.TargetAgent = a
		} else {
			gr.sendError(execCtx, "no_available_agent", "No agent configured")
			return
		}
		go gr.core.Stream(ctx, fallback)
		return
	}
	go gr.orchestrateConversation(ctx, execCtx, mainSession, availableMembers)
}

func (gr *GroupRouter) orchestrateConversation(ctx context.Context, execCtx ExecutionContext, mainSession *Session, availableMembers []GroupMember) {
	groupID := execCtx.Session.GroupInfo.GroupID

	// Track the current message context flowing through the chain
	currentMessageContext := fmt.Sprintf("【用户初始问题】：\n%s\n", execCtx.Message.Content.Text)

	// We will loop up to 5 times to prevent infinite loops (configurable)
	maxIterations := 5
	for i := 0; i < maxIterations; i++ {
		// Build Prompt for LLM routing
		routingPrompt := gr.buildRoutingPrompt(execCtx, availableMembers, currentMessageContext, i > 0)

		// Route Context
		routeCtx := execCtx
		routeCtx.Session = *mainSession
		routeCtx.Message.Content.Text = routingPrompt
		routeCtx.Message.Metadata.OverrideConfig = &protocol.UserInteractionConfig{
			StreamMode:     "final_only",
			EnableThinking: false, // Don't stream thinking for routing
		}
		routeCtx.SystemPrompt = "你是一个群聊调度总管(MainAgent)。你要根据群成员的能力，规划问题的解答路径。只需返回JSON格式的【一个】成员ID数组作为下一步执行者，或者空数组[]表示所有任务已完成。无需任何解释！"

		// Select main agent (default logic)
		mainAgentConfig, err := gr.selectAgent(routeCtx)
		if err != nil {
			gr.sendError(execCtx, "no_available_agent", err.Error())
			return
		}
		routeCtx.TargetAgent = mainAgentConfig

		// Create subcontext
		subCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		log.Printf("[group_router] Routing iteration %d for group=%s msg=%s", i+1, groupID, execCtx.Message.MessageID)
		// Execute internal stream
		final, _, upstreamSessionID, err := gr.core.streamWithSessionRecovery(subCtx, routeCtx, streamSendNone, "ROUTING", false, true)
		cancel()

		if err != nil {
			log.Printf("[group_router] MainAgent routing failed: %v", err)
			break
		}
		if strings.TrimSpace(upstreamSessionID) != "" && mainSession.UpstreamSessionID != upstreamSessionID {
			mainSession.UpstreamSessionID = upstreamSessionID
			_ = gr.core.sessionMgr.UpdateSession(mainSession)
		}

		// Parse picked IDs
		pickedIDs := gr.parseRoutingResult(final, availableMembers)
		log.Printf("[group_router] MainAgent selected members: %v from final: %q", pickedIDs, final)

		if len(pickedIDs) == 0 {
			// Finished
			log.Printf("[group_router] MainAgent returned empty array, ending chain.")
			break
		}

		// Take the first picked agent to act in this turn (expert chain)
		targetUserID := pickedIDs[0]
		var targetMember *GroupMember
		for _, m := range availableMembers {
			if m.UserID == targetUserID {
				mCopy := m
				targetMember = &mCopy
				break
			}
		}

		if targetMember == nil {
			log.Printf("[group_router] Picked invalid member %s, stopping chain", targetUserID)
			break
		}

		// Dispatch and wait for response
		reply, err := gr.dispatchToMemberAndWait(ctx, execCtx, *targetMember, availableMembers, currentMessageContext)
		if err != nil {
			log.Printf("[group_router] Member execution failed: %v", err)
			break
		}

		// Append to chain context
		currentMessageContext += fmt.Sprintf("\n【%s(群成员) 的执行结果】：\n%s\n", targetMember.UserID, reply)
	}
	log.Printf("[group_router] Group orchestration finished for msg=%s", execCtx.Message.MessageID)
}

// Executes an agent's turn blockingly and returns the final markdown response string
func (gr *GroupRouter) dispatchToMemberAndWait(ctx context.Context, execCtx ExecutionContext, member GroupMember, allMembers []GroupMember, currentContext string) (string, error) {
	memCtx := gr.buildMemberContext(execCtx, member, allMembers)

	agentConfig, err := gr.selectAgent(memCtx)
	if err != nil {
		return "", err
	}
	memCtx.TargetAgent = agentConfig

	// Create a new message containing the full chain context so far
	memCtx.Message.Content.Text = fmt.Sprintf("你正在协同处理一个群组任务，以下是目前的进展情况和初始问题，你需要接着往下处理，并直接给出你的专注领域的解题步骤：\n\n%s", currentContext)

	// Create an interceptor emitter so we can stream to the user AND capture the final string
	originalEmitter := gr.core.emitter
	wrappedEmitter := &groupMemberEmitter{
		original:       originalEmitter,
		groupSessionID: execCtx.Session.SessionID,
		senderID:       member.UserID,
		buffer:         &strings.Builder{},
	}

	gr.core.SetEmitter(wrappedEmitter)

	// Force the AgentCore to stream back SSE (and to our interceptor)
	log.Printf("[group_router] Handing over execution to member %s for message %s", member.UserID, memCtx.Message.MessageID)

	subCtx, cancel := context.WithTimeout(ctx, 300*time.Second) // Long timeout for complex agent tasks
	defer cancel()

	// We call streamOnce natively because we want the final response to feed into next iteration
	mode := streamSendRealtime
	if memCtx.EffectiveConfig.StreamMode == "final_only" {
		mode = streamSendFinalOnly
	}

	finalResponse, _, upstreamSessionID, err := gr.core.streamWithSessionRecovery(subCtx, memCtx, mode, "GROUP_EXPERT_FINAL", true, false)

	gr.core.SetEmitter(originalEmitter)

	if strings.TrimSpace(upstreamSessionID) != "" {
		sid := gr.store.GetMemberSession(execCtx.Session.GroupInfo.GroupID, member.UserID)
		if sess, err := gr.core.sessionMgr.GetOrCreateSession(member.UserID, sid, "single"); err == nil && sess != nil && sess.UpstreamSessionID != upstreamSessionID {
			sess.UpstreamSessionID = upstreamSessionID
			_ = gr.core.sessionMgr.UpdateSession(sess)
		}
	}

	return finalResponse, err
}

type groupMemberEmitter struct {
	original       Emitter
	groupSessionID string
	senderID       string
	buffer         *strings.Builder
}

func (e *groupMemberEmitter) SendToConnection(connID string, msg protocol.ServerMessage) {
	e.original.BroadcastToSession(e.groupSessionID, e.injectSenderID(msg))
}

func (e *groupMemberEmitter) BroadcastToSession(sessionID string, msg protocol.ServerMessage) {
	e.original.BroadcastToSession(e.groupSessionID, e.injectSenderID(msg))
}

func (e *groupMemberEmitter) injectSenderID(msg protocol.ServerMessage) protocol.ServerMessage {
	msg.SessionID = e.groupSessionID
	if payload, ok := msg.Payload.(protocol.StreamChunkPayload); ok {
		payload.SenderID = e.senderID
		msg.Payload = payload
	} else if payload, ok := msg.Payload.(protocol.StreamEndPayload); ok {
		payload.SenderID = e.senderID
		e.buffer.WriteString(payload.Content) // capture the final content
		msg.Payload = payload
	}
	return msg
}

func (gr *GroupRouter) buildMemberContext(ctx ExecutionContext, member GroupMember, allMembers []GroupMember) ExecutionContext {
	groupID := ctx.Session.GroupInfo.GroupID

	// Ensure member session
	memberSessionID := gr.store.GetMemberSession(groupID, member.UserID)
	memberSession, err := gr.core.sessionMgr.GetOrCreateSession(member.UserID, memberSessionID, "single")
	if err == nil {
		gr.store.SetMemberSession(groupID, member.UserID, memberSession.SessionID)
	} else {
		memberSession = &Session{SessionID: memberSessionID, UserID: member.UserID, Type: "single"}
	}

	newCtx := ctx
	newCtx.Session = *memberSession

	// Create identity
	newCtx.RoleIdentity = &RoleIdentity{
		UserID:          member.UserID,
		RoleDescription: member.Description,
		Skills:          member.Skills,
		Persona:         member.Persona,
	}

	// Fetch backend config for this agent
	beAgent, err := gr.core.fetchBackendAgent(context.Background(), member.UserID)
	if err == nil && beAgent != nil {
		newCtx.AgentProfile = buildAgentSnapshot(*beAgent)
		newCtx.UserProfile.ModelID = beAgent.Model
		newCtx.UserProfile.SimplifiedOutput = beAgent.SimplifiedOutput
		newCtx.EffectiveConfig.EnableThinking = beAgent.ThinkingEnabled
		if newCtx.Message.Metadata.ModelID == "" && beAgent.Model != "" {
			newCtx.Message.Metadata.ModelID = beAgent.Model
		}
	}

	// Build a super prompt that includes other members' descriptions
	var skillDescs []string
	for _, s := range member.Skills {
		skillDescs = append(skillDescs, fmt.Sprintf("- %s: %s", s.SkillName, s.Description))
	}

	var otherMembersStr strings.Builder
	for _, m := range allMembers {
		if m.UserID != member.UserID {
			otherMembersStr.WriteString(fmt.Sprintf("- [%s]: %s\n", m.UserID, m.Description))
		}
	}

	agentPersona := ""
	if beAgent != nil {
		agentPersona = buildAgentPersonaPrompt(*beAgent)
	} else {
		agentPersona = fmt.Sprintf("你是群聊核心成员 \"%s\"\n职责：%s\n", member.UserID, member.Description)
	}

	newCtx.UserProfile.SystemPrompt = fmt.Sprintf(`%s

# 团队协作环境说明
你正在一个协助网络(Agent Team)内工作。群内有其他专家成员，你可以放心地在输出结尾，说明下一步需要哪个成员来接手。
其他团队成员及能力简介：
%s

# 执行约定
严格基于你的职责领域进行解答。如果有需要其他成员协作的中间步骤，可以在执行完你的专属部分后，给出详细的交接结果返回。`,
		agentPersona,
		otherMembersStr.String(),
	)

	return newCtx
}

func (gr *GroupRouter) buildRoutingPrompt(ctx ExecutionContext, members []GroupMember, currentContext string, isHandover bool) string {
	mode := "自由回答模式(free)"
	if ctx.Session.GroupInfo.GroupRuleMode != "" {
		mode = ctx.Session.GroupInfo.GroupRuleMode
	}

	var membersJson []map[string]interface{}
	for _, m := range members {
		membersJson = append(membersJson, map[string]interface{}{
			"id":          m.UserID,
			"description": m.Description,
		})
	}
	mbBytes, _ := json.MarshalIndent(membersJson, "", "  ")

	actionText := ""
	if isHandover {
		actionText = "当前任务仍在处理中，已经得到部分成员的回复。请阅读上下文，判断是否有下一步工作。若需要交给下一个专家成员继续，返回由该成员ID构成的数组，比如 [\"user2\"]。如果任务已彻底完结且不需要继续讨论，请返回空数组 [] 停止！"
	} else {
		actionText = "用户发送了新问题，请根据成员画像分配第一个专家成员来处理。返回包含被选员ID的数组，比如 [\"user1\"]"
	}

	prompt := fmt.Sprintf(`你是群组的总管调度引擎。请仔细阅读当前的协作进展和群规则，并分发任务。

【群聊规则模式】：%s
- 如果是 free (自由回答)，选择最能推进解答的人。
- 如果是 expert (专家模式)，选择唯一的权威专家来作答。
- 如果是 custom (自定义模式)，遵循自定义规则：%s

【可用的群成员画像库】：
%s

【当前协同工作流/接力状态】：
%s

【决策指引】：
%s
请严格输出唯一一个合法的JSON字符串数组，比如 ["userX"]。若已解决，输出 []。不要包含任何 Markdown 代码块标记或其他解释！`,
		mode,
		ctx.Session.GroupInfo.CustomRule,
		string(mbBytes),
		currentContext,
		actionText)

	return prompt
}

func (gr *GroupRouter) parseRoutingResult(result string, available []GroupMember) []string {
	result = strings.TrimSpace(result)
	start := strings.Index(result, "[")
	end := strings.LastIndex(result, "]")
	if start >= 0 && end > start {
		result = result[start : end+1]
	} else {
		// Try to parse raw token as single element array if it matches an id
		for _, m := range available {
			if strings.Contains(result, m.UserID) {
				return []string{m.UserID}
			}
		}
		return nil
	}

	var rawIDs []string
	if err := json.Unmarshal([]byte(result), &rawIDs); err != nil {
		log.Printf("[group_router] JSON parse failed for routing result %q: %v", result, err)
		return nil
	}

	// Validate IDs
	var valid []string
	for _, id := range rawIDs {
		id = strings.TrimSpace(id)
		for _, m := range available {
			if m.UserID == id {
				valid = append(valid, id)
				break
			}
		}
	}
	return valid
}

func (gr *GroupRouter) selectAgent(ctx ExecutionContext) (SubAgentConfig, error) {
	if ctx.Message.Metadata.TargetAgentID != "" {
		a, ok := gr.core.pool.Get(ctx.Message.Metadata.TargetAgentID)
		if ok {
			return a, nil
		}
	}
	a, ok := gr.core.pool.Default()
	if !ok {
		return SubAgentConfig{}, errors.New("no agent configured")
	}
	return a, nil
}

func (gr *GroupRouter) sendError(execCtx ExecutionContext, code, msg string) {
	if gr.core.emitter != nil {
		gr.core.emitter.SendToConnection(execCtx.ConnectionID, protocol.ServerMessage{
			Type:      "error",
			MessageID: execCtx.Message.MessageID,
			SessionID: execCtx.Session.SessionID,
			Timestamp: time.Now().UnixMilli(),
			Payload: protocol.ErrorPayload{
				Code:      code,
				Message:   msg,
				Retryable: false,
			},
		})
	}
}
