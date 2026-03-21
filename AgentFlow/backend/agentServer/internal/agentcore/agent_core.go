package agentcore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"agentServer/internal/opencode"
	"agentServer/internal/protocol"
)

type Emitter interface {
	SendToConnection(connID string, msg protocol.ServerMessage)
	BroadcastToSession(sessionID string, msg protocol.ServerMessage)
}

type AgentCore struct {
	pool             *SubAgentPool
	streamManager    *StreamManager
	messageRouter    *MessageRouter
	sessionMgr       *SessionManager
	configMgr        *ConfigManager
	logPersistence   *LogPersistence
	groupRouter      *GroupRouter
	singleChat       *SingleChatHandler
	backendModelsURL string
	emitter          Emitter

	clientsMu  sync.RWMutex
	clients    map[string]*opencode.Client
	httpClient *http.Client
}

func NewAgentCore(agents []SubAgentConfig, backendModelsURL string, httpClient *http.Client) *AgentCore {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	pool := NewSubAgentPool(agents)
	ac := &AgentCore{
		pool:             pool,
		streamManager:    NewStreamManager(),
		sessionMgr:       NewSessionManager(),
		configMgr:        NewConfigManager(),
		logPersistence:   NewLogPersistence(),
		backendModelsURL: backendModelsURL,
		clients:          make(map[string]*opencode.Client),
		httpClient:       httpClient,
	}
	ac.singleChat = NewSingleChatHandler(ac)
	ac.groupRouter = NewGroupRouter(ac)
	ac.messageRouter = NewMessageRouter(ac.singleChat, ac.groupRouter)
	return ac
}

func (ac *AgentCore) SetEmitter(emitter Emitter) {
	ac.emitter = emitter
}

func (ac *AgentCore) SetDefaultUserConfig(cfg UserInteractionConfig) {
	ac.configMgr.SetDefault(cfg)
}

func (ac *AgentCore) SetLogBasePath(basePath string) {
	ac.logPersistence.SetBasePath(basePath)
}

func (ac *AgentCore) HandleMessage(ctx context.Context, connID string, msg protocol.Message) {
	log.Printf("[agentcore] handle message conn_id=%s message_id=%s session_id=%s user_id=%s group_id=%s type=%s", connID, msg.MessageID, msg.SessionID, msg.UserID, msg.GroupID, msg.Type)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[agentcore] panic conn_id=%s message_id=%s err=%v", connID, msg.MessageID, r)
			ac.sendError(connID, msg.SessionID, msg.MessageID, "internal_error", fmt.Sprintf("%v", r), false)
		}
	}()

	sessionType := detectMsgType(msg)
	session, err := ac.sessionMgr.GetOrCreateSession(msg.UserID, msg.SessionID, sessionType)
	if err != nil {
		log.Printf("[agentcore] session error conn_id=%s message_id=%s err=%v", connID, msg.MessageID, err)
		ac.sendError(connID, msg.SessionID, msg.MessageID, "session_error", err.Error(), false)
		if msg.UserID != "" {
			_ = ac.logPersistence.AppendLog(msg.UserID, msg.SessionID, "ERROR", err.Error())
		}
		return
	}

	userProfile, err := ac.configMgr.GetUserProfile(msg.UserID)
	if err != nil {
		userProfile = ac.configMgr.GetDefaultProfile(msg.UserID)
	}

	if sessionType == "single" && msg.UserID != "" {
		if backendAgent, err := ac.fetchBackendAgent(ctx, msg.UserID); err == nil && backendAgent != nil {
			userProfile.SystemPrompt = backendAgent.SystemPrompt
			userProfile.ModelID = backendAgent.Model
			userProfile.SimplifiedOutput = backendAgent.SimplifiedOutput
			userProfile.Config.EnableThinking = backendAgent.ThinkingEnabled

			if msg.Metadata.ModelID == "" && backendAgent.Model != "" {
				msg.Metadata.ModelID = backendAgent.Model
			}
			log.Printf("[agentcore] loaded backend agent config user_id=%s model_id=%s thinking=%t simplified=%t system_prompt=%t", msg.UserID, msg.Metadata.ModelID, backendAgent.ThinkingEnabled, backendAgent.SimplifiedOutput, strings.TrimSpace(backendAgent.SystemPrompt) != "")
		} else if err != nil {
			log.Printf("[agentcore] backend agent config load failed user_id=%s err=%v", msg.UserID, err)
		}
	}

	if userProfile.CurrentSessionID != session.SessionID {
		userProfile.CurrentSessionID = session.SessionID
	}
	_ = ac.configMgr.UpdateUserProfile(userProfile)

	effective := mergeUserConfig(ac.configMgr.DefaultConfig(), userProfile.Config, msg.Metadata.OverrideConfig)
	systemPrompt := ""
	if sessionType == "single" {
		systemPrompt = buildSingleSystemPrompt(userProfile.SystemPrompt, effective.EnableThinking)
	}
	execCtx := ExecutionContext{
		Message:         msg,
		Session:         *session,
		UserProfile:     *userProfile,
		ConnectionID:    connID,
		SystemPrompt:    systemPrompt,
		EffectiveConfig: effective,
	}
	ac.logExecutionSnapshot(execCtx)

	ac.messageRouter.Route(ctx, execCtx)
	log.Printf("[agentcore] route scheduled conn_id=%s message_id=%s session_id=%s session_type=%s", connID, msg.MessageID, execCtx.Session.SessionID, sessionType)
}

func (ac *AgentCore) getClient(endpoint string) *opencode.Client {
	ac.clientsMu.RLock()
	c, ok := ac.clients[endpoint]
	ac.clientsMu.RUnlock()
	if ok {
		return c
	}

	ac.clientsMu.Lock()
	defer ac.clientsMu.Unlock()
	if c, ok := ac.clients[endpoint]; ok {
		return c
	}
	c = opencode.NewClient(endpoint, ac.backendModelsURL, ac.httpClient)
	ac.clients[endpoint] = c
	return c
}

func (ac *AgentCore) Stream(ctx context.Context, execCtx ExecutionContext) {
	if ac.emitter == nil {
		log.Printf("[agentcore] stream skipped: emitter missing conn_id=%s message_id=%s", execCtx.ConnectionID, execCtx.Message.MessageID)
		return
	}
	log.Printf("[agentcore] stream start conn_id=%s session_id=%s message_id=%s agent=%s return_mode=%s stream_mode=%s", execCtx.ConnectionID, execCtx.Session.SessionID, execCtx.Message.MessageID, execCtx.TargetAgent.AgentID, execCtx.EffectiveConfig.ReturnStrategy.Type, execCtx.EffectiveConfig.StreamMode)

	_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "USER", execCtx.Message.Content.Text)
	if execCtx.SystemPrompt != "" {
		_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "SYSTEM", execCtx.SystemPrompt)
	}

	returnMode := execCtx.EffectiveConfig.ReturnStrategy.Type
	if returnMode == "" {
		returnMode = execCtx.EffectiveConfig.StreamMode
	}
	allowChunks := execCtx.EffectiveConfig.StreamMode == "realtime" && returnMode != "final_only"

	send := func(m protocol.ServerMessage) {
		if execCtx.Session.Type == "group" {
			ac.emitter.BroadcastToSession(execCtx.Session.SessionID, m)
			return
		}
		ac.emitter.SendToConnection(execCtx.ConnectionID, m)
	}

	if execCtx.Session.Type == "single" && execCtx.UserProfile.SimplifiedOutput {
		rawFinal, err := ac.streamOnce(ctx, execCtx, streamSendNone, "RAW_FINAL", false)
		if err != nil {
			log.Printf("[agentcore] simplified raw stream error conn_id=%s message_id=%s err=%v", execCtx.ConnectionID, execCtx.Message.MessageID, err)
			send(protocol.ServerMessage{
				Type:      "error",
				MessageID: execCtx.Message.MessageID,
				SessionID: execCtx.Session.SessionID,
				Payload: protocol.ErrorPayload{
					Code:      "opencode_error",
					Message:   err.Error(),
					Retryable: true,
				},
			})
			return
		}

		summaryMsg := execCtx.Message
		summaryMsg.Content.Text = buildSummaryRequest(execCtx.Message.Content.Text, rawFinal)
		_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "USER_SUMMARY", summaryMsg.Content.Text)

		summaryCtx := execCtx
		summaryCtx.Message = summaryMsg
		summaryCtx.UserProfile.SimplifiedOutput = false
		summaryCtx.EffectiveConfig.EnableThinking = false
		summaryCtx.SystemPrompt = buildSingleSystemPrompt(execCtx.UserProfile.SystemPrompt, false)

		summaryFinal, err := ac.streamOnce(ctx, summaryCtx, streamSendFinalOnly, "SUMMARY", true)
		if err != nil {
			log.Printf("[agentcore] simplified summary error conn_id=%s message_id=%s err=%v", execCtx.ConnectionID, execCtx.Message.MessageID, err)
			_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "ERROR", err.Error())
			send(protocol.ServerMessage{
				Type:      "error",
				MessageID: execCtx.Message.MessageID,
				SessionID: execCtx.Session.SessionID,
				Payload: protocol.ErrorPayload{
					Code:      "opencode_error",
					Message:   err.Error(),
					Retryable: true,
				},
			})
			return
		}

		send(protocol.ServerMessage{
			Type:      "stream_end",
			MessageID: execCtx.Message.MessageID,
			SessionID: execCtx.Session.SessionID,
			Timestamp: time.Now().UnixMilli(),
			Payload: protocol.StreamEndPayload{
				Content: summaryFinal,
				IsFinal: true,
			},
		})
		return
	}

	mode := streamSendFinalOnly
	if allowChunks {
		mode = streamSendRealtime
	}
	final, err := ac.streamOnce(ctx, execCtx, mode, "FINAL", true)
	if err != nil {
		log.Printf("[agentcore] stream error conn_id=%s message_id=%s err=%v", execCtx.ConnectionID, execCtx.Message.MessageID, err)
		_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "ERROR", err.Error())
		send(protocol.ServerMessage{
			Type:      "error",
			MessageID: execCtx.Message.MessageID,
			SessionID: execCtx.Session.SessionID,
			Payload: protocol.ErrorPayload{
				Code:      "opencode_error",
				Message:   err.Error(),
				Retryable: true,
			},
		})
		return
	}

	send(protocol.ServerMessage{
		Type:      "stream_end",
		MessageID: execCtx.Message.MessageID,
		SessionID: execCtx.Session.SessionID,
		Timestamp: time.Now().UnixMilli(),
		Payload: protocol.StreamEndPayload{
			Content: final,
			IsFinal: true,
		},
	})
	log.Printf("[agentcore] stream end sent conn_id=%s message_id=%s session_id=%s final_len=%d", execCtx.ConnectionID, execCtx.Message.MessageID, execCtx.Session.SessionID, len(final))
}

func logType(isThinking bool, isFinal bool) string {
	if isFinal {
		return "FINAL"
	}
	if isThinking {
		return "THINKING"
	}
	return "CHUNK"
}

func (ac *AgentCore) sendError(connID string, sessionID string, messageID string, code string, message string, retryable bool) {
	if ac.emitter == nil {
		return
	}
	log.Printf("[agentcore] send error conn_id=%s session_id=%s message_id=%s code=%s retryable=%t message=%s", connID, sessionID, messageID, code, retryable, message)
	ac.emitter.SendToConnection(connID, protocol.ServerMessage{
		Type:      "error",
		SessionID: sessionID,
		MessageID: messageID,
		Payload: protocol.ErrorPayload{
			Code:      code,
			Message:   message,
			Retryable: retryable,
		},
	})
}

func (ac *AgentCore) logExecutionSnapshot(execCtx ExecutionContext) {
	systemPromptPresent := strings.TrimSpace(execCtx.UserProfile.SystemPrompt) != ""
	content := fmt.Sprintf(
		"user_id=%s session_id=%s message_id=%s model_id=%s thinking=%t simplified=%t stream_mode=%s return_type=%s system_prompt=%t",
		execCtx.UserProfile.UserID,
		execCtx.Session.SessionID,
		execCtx.Message.MessageID,
		execCtx.UserProfile.ModelID,
		execCtx.EffectiveConfig.EnableThinking,
		execCtx.UserProfile.SimplifiedOutput,
		execCtx.EffectiveConfig.StreamMode,
		execCtx.EffectiveConfig.ReturnStrategy.Type,
		systemPromptPresent,
	)
	_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "CONFIG", content)
}

func detectMsgType(msg protocol.Message) string {
	if msg.GroupID != "" {
		return "group"
	}
	if len(msg.Content.Mentions) > 1 {
		return "group"
	}
	if msg.Type == "group" {
		return "group"
	}
	return "single"
}

type backendAgent struct {
	ID               string `json:"id"`
	Model            string `json:"model"`
	SystemPrompt     string `json:"system_prompt"`
	ThinkingEnabled  bool   `json:"thinking_enabled"`
	SimplifiedOutput bool   `json:"simplified_output"`
}

func (ac *AgentCore) fetchBackendAgent(ctx context.Context, userID string) (*backendAgent, error) {
	base := backendBaseURL(ac.backendModelsURL)
	if base == "" {
		return nil, fmt.Errorf("backend url not configured")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/api/agents/", nil)
	if err != nil {
		return nil, err
	}
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.ReadAll(resp.Body)
		return nil, fmt.Errorf("backend agents error (%d)", resp.StatusCode)
	}

	var agents []backendAgent
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		return nil, err
	}
	for i := range agents {
		if agents[i].ID == userID {
			return &agents[i], nil
		}
	}
	return nil, fmt.Errorf("agent not found")
}

func backendBaseURL(modelsAPIURL string) string {
	modelsAPIURL = strings.TrimSpace(modelsAPIURL)
	if modelsAPIURL == "" {
		return ""
	}
	u, err := url.Parse(modelsAPIURL)
	if err != nil {
		return ""
	}
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""
	return strings.TrimRight(u.String(), "/")
}

func buildSingleSystemPrompt(base string, enableThinking bool) string {
	base = strings.TrimSpace(base)
	if base == "" {
		if enableThinking {
			return "请在回答中使用 <think>...</think> 包裹思考过程，然后输出最终答案。"
		}
		return "请直接输出最终答案，不要输出思考过程。"
	}
	if enableThinking {
		return base + "\n\n请在回答中使用 <think>...</think> 包裹思考过程，然后输出最终答案。"
	}
	return base + "\n\n请直接输出最终答案，不要输出思考过程。"
}

func buildSummaryRequest(userText string, rawFinal string) string {
	userText = strings.TrimSpace(userText)
	rawFinal = strings.TrimSpace(rawFinal)
	if userText == "" && rawFinal == "" {
		return "请输出最终结果摘要，不要输出思考过程。"
	}
	return "请基于下面的内容，给出简洁的最终结果摘要（不要输出思考过程），只输出答案：\n\n[用户问题]\n" + userText + "\n\n[完整回答]\n" + rawFinal
}

type streamSendMode int

const (
	streamSendNone streamSendMode = iota
	streamSendRealtime
	streamSendFinalOnly
)

type chunkFragment struct {
	content    string
	isThinking bool
}

type thinkSplitter struct {
	inThinking bool
	carry      string
}

func (p *thinkSplitter) Split(chunk string) []chunkFragment {
	const startTag = "<think>"
	const endTag = "</think>"

	combined := p.carry + chunk
	p.carry = ""

	carryLen := suffixCarryLen(combined, startTag, endTag)
	if carryLen > 0 {
		p.carry = combined[len(combined)-carryLen:]
		combined = combined[:len(combined)-carryLen]
	}

	var out []chunkFragment
	for len(combined) > 0 {
		if !p.inThinking {
			i := strings.Index(combined, startTag)
			if i < 0 {
				out = append(out, chunkFragment{content: combined, isThinking: false})
				return out
			}
			if i > 0 {
				out = append(out, chunkFragment{content: combined[:i], isThinking: false})
			}
			combined = combined[i+len(startTag):]
			p.inThinking = true
			continue
		}

		i := strings.Index(combined, endTag)
		if i < 0 {
			out = append(out, chunkFragment{content: combined, isThinking: true})
			return out
		}
		if i > 0 {
			out = append(out, chunkFragment{content: combined[:i], isThinking: true})
		}
		combined = combined[i+len(endTag):]
		p.inThinking = false
	}

	return out
}

func suffixCarryLen(s string, tokens ...string) int {
	max := 0
	for _, tok := range tokens {
		if tok == "" {
			continue
		}
		for i := 1; i < len(tok); i++ {
			if strings.HasSuffix(s, tok[:i]) {
				if i > max {
					max = i
				}
			}
		}
	}
	return max
}

func (ac *AgentCore) streamOnce(ctx context.Context, execCtx ExecutionContext, mode streamSendMode, finalLogType string, finalize bool) (string, error) {
	reqMsg := execCtx.Message
	if execCtx.SystemPrompt != "" {
		reqMsg.Content.Text = execCtx.SystemPrompt + "\n\n" + reqMsg.Content.Text
	}

	client := ac.getClient(execCtx.TargetAgent.Endpoint)
	streamCh, errCh := client.Stream(ctx, reqMsg)
	log.Printf("[agentcore] upstream stream opened conn_id=%s message_id=%s endpoint=%s", execCtx.ConnectionID, execCtx.Message.MessageID, execCtx.TargetAgent.Endpoint)

	send := func(m protocol.ServerMessage) {
		if mode == streamSendNone {
			return
		}
		if execCtx.Session.Type == "group" {
			ac.emitter.BroadcastToSession(execCtx.Session.SessionID, m)
			return
		}
		ac.emitter.SendToConnection(execCtx.ConnectionID, m)
	}

	splitter := &thinkSplitter{}
	var visible strings.Builder
	var streamOpen = true
	var errOpen = true

	for streamOpen || errOpen {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case err, ok := <-errCh:
			if !ok {
				errOpen = false
				continue
			}
			if err != nil {
				log.Printf("[agentcore] upstream error conn_id=%s message_id=%s err=%v", execCtx.ConnectionID, execCtx.Message.MessageID, err)
				return "", err
			}
		case ev, ok := <-streamCh:
			if !ok {
				streamOpen = false
				continue
			}

			switch ev.Type {
			case "stream_chunk":
				payload, _ := ev.Payload.(protocol.StreamChunkPayload)
				log.Printf("[agentcore] upstream chunk conn_id=%s message_id=%s len=%d is_thinking=%t", execCtx.ConnectionID, execCtx.Message.MessageID, len(payload.Content), payload.IsThinking)
				var frags []chunkFragment
				if payload.IsThinking {
					frags = []chunkFragment{{content: payload.Content, isThinking: true}}
				} else {
					frags = splitter.Split(payload.Content)
				}
				for _, frag := range frags {
					if frag.content == "" {
						continue
					}
					if frag.isThinking {
						_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "THINKING", frag.content)
					} else {
						visible.WriteString(frag.content)
						_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "CHUNK", frag.content)
					}

					if mode == streamSendRealtime && (!frag.isThinking || execCtx.EffectiveConfig.EnableThinking) {
						send(protocol.ServerMessage{
							Type:      "stream_chunk",
							MessageID: execCtx.Message.MessageID,
							SessionID: execCtx.Session.SessionID,
							Timestamp: time.Now().UnixMilli(),
							Payload: protocol.StreamChunkPayload{
								Content:    frag.content,
								IsThinking: frag.isThinking,
								IsFinal:    false,
								Metadata:   payload.Metadata,
							},
						})
					}
				}
			case "stream_end":
				final := visible.String()
				log.Printf("[agentcore] upstream end conn_id=%s message_id=%s final_len=%d", execCtx.ConnectionID, execCtx.Message.MessageID, len(final))
				if finalLogType != "" {
					_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, finalLogType, final)
				}
				if finalize {
					_ = ac.logPersistence.FinalizeLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, len(final))
				}
				return final, nil
			case "error":
				payload, _ := ev.Payload.(protocol.ErrorPayload)
				log.Printf("[agentcore] upstream returned error conn_id=%s message_id=%s code=%s message=%s", execCtx.ConnectionID, execCtx.Message.MessageID, payload.Code, payload.Message)
				return "", fmt.Errorf("%s", payload.Message)
			}
		}
	}

	final := visible.String()
	if finalLogType != "" {
		_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, finalLogType, final)
	}
	if finalize {
		_ = ac.logPersistence.FinalizeLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, len(final))
	}
	return final, nil
}
