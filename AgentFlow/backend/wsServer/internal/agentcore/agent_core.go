package agentcore

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"wsServer/internal/opencode"
	"wsServer/internal/protocol"
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
	defer func() {
		if r := recover(); r != nil {
			ac.sendError(connID, msg.SessionID, msg.MessageID, "internal_error", fmt.Sprintf("%v", r), false)
		}
	}()

	sessionType := detectMsgType(msg)
	session, err := ac.sessionMgr.GetOrCreateSession(msg.UserID, msg.SessionID, sessionType)
	if err != nil {
		ac.sendError(connID, msg.SessionID, msg.MessageID, "session_error", err.Error(), false)
		return
	}

	userProfile, err := ac.configMgr.GetUserProfile(msg.UserID)
	if err != nil {
		userProfile = ac.configMgr.GetDefaultProfile(msg.UserID)
	}

	if userProfile.CurrentSessionID != session.SessionID {
		userProfile.CurrentSessionID = session.SessionID
		_ = ac.configMgr.UpdateUserProfile(userProfile)
	}

	effective := mergeUserConfig(ac.configMgr.DefaultConfig(), userProfile.Config, msg.Metadata.OverrideConfig)
	execCtx := ExecutionContext{
		Message:         msg,
		Session:         *session,
		UserProfile:     *userProfile,
		ConnectionID:    connID,
		EffectiveConfig: effective,
	}

	ac.messageRouter.Route(ctx, execCtx)
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
		return
	}

	reqMsg := execCtx.Message
	if execCtx.SystemPrompt != "" {
		reqMsg.Content.Text = execCtx.SystemPrompt + "\n\n" + reqMsg.Content.Text
	}

	client := ac.getClient(execCtx.TargetAgent.Endpoint)
	streamCh, errCh := client.Stream(ctx, reqMsg)
	var streamOpen = true
	var errOpen = true
	var total strings.Builder
	returnMode := execCtx.EffectiveConfig.ReturnStrategy.Type
	if returnMode == "" {
		returnMode = execCtx.EffectiveConfig.StreamMode
	}

	for streamOpen || errOpen {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-errCh:
			if !ok {
				errOpen = false
				continue
			}
			if err != nil {
				ac.emitter.SendToConnection(execCtx.ConnectionID, protocol.ServerMessage{
					Type:      "error",
					MessageID: execCtx.Message.MessageID,
					SessionID: execCtx.Session.SessionID,
					Payload: protocol.ErrorPayload{
						Code:      "opencode_error",
						Message:   err.Error(),
						Retryable: true,
					},
				})
			}
			return
		case ev, ok := <-streamCh:
			if !ok {
				streamOpen = false
				continue
			}

			send := func(m protocol.ServerMessage) {
				if execCtx.Session.Type == "group" {
					ac.emitter.BroadcastToSession(execCtx.Session.SessionID, m)
					return
				}
				ac.emitter.SendToConnection(execCtx.ConnectionID, m)
			}

			switch ev.Type {
			case "stream_chunk":
				payload, _ := ev.Payload.(protocol.StreamChunkPayload)
				total.WriteString(payload.Content)
				_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, logType(payload.IsThinking, payload.IsFinal), payload.Content)

				if execCtx.EffectiveConfig.StreamMode == "realtime" && returnMode != "final_only" {
					send(ev)
				}
			case "stream_end":
				endPayload, _ := ev.Payload.(protocol.StreamEndPayload)
				if endPayload.Content == "" {
					endPayload.Content = total.String()
				}
				_ = ac.logPersistence.AppendLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, "FINAL", endPayload.Content)
				_ = ac.logPersistence.FinalizeLog(execCtx.UserProfile.UserID, execCtx.Session.SessionID, len(endPayload.Content))

				outEv := ev
				outEv.Payload = endPayload
				outEv.Timestamp = time.Now().UnixMilli()
				send(outEv)
				return
			case "error":
				send(ev)
				return
			default:
				send(ev)
			}
		}
	}
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

func detectMsgType(msg protocol.Message) string {
	if len(msg.Content.Mentions) > 1 {
		return "group"
	}
	if msg.Type == "group" {
		return "group"
	}
	return "single"
}
