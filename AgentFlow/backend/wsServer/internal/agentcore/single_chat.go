package agentcore

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"wsServer/internal/protocol"
)

type SingleChatHandler struct {
	core *AgentCore
}

func NewSingleChatHandler(core *AgentCore) *SingleChatHandler {
	return &SingleChatHandler{core: core}
}

func (h *SingleChatHandler) Handle(ctx context.Context, execCtx ExecutionContext) {
	agent, err := h.selectAgent(execCtx)
	if err != nil {
		if h.core.emitter != nil {
			h.core.emitter.SendToConnection(execCtx.ConnectionID, protocol.ServerMessage{
				Type:      "error",
				MessageID: execCtx.Message.MessageID,
				SessionID: execCtx.Session.SessionID,
				Timestamp: time.Now().UnixMilli(),
				Payload: protocol.ErrorPayload{
					Code:      "no_available_agent",
					Message:   err.Error(),
					Retryable: false,
				},
			})
		}
		return
	}

	execCtx.TargetAgent = agent
	execCtx.Session.LastAgentID = agent.AgentID
	_ = h.core.sessionMgr.UpdateSession(&execCtx.Session)

	h.core.Stream(ctx, execCtx)
}

func (h *SingleChatHandler) selectAgent(ctx ExecutionContext) (SubAgentConfig, error) {
	if ctx.Message.Metadata.TargetAgentID != "" {
		agent, ok := h.core.pool.Get(ctx.Message.Metadata.TargetAgentID)
		if ok {
			return agent, nil
		}
	}

	text := ctx.Message.Content.Text
	agents := h.core.pool.GetAll()
	if len(agents) == 0 {
		return SubAgentConfig{}, errors.New("no agent configured")
	}

	type scored struct {
		agent SubAgentConfig
		score float64
	}

	out := make([]scored, 0, len(agents))
	for _, a := range agents {
		score := 0.0

		for _, cap := range a.Capabilities {
			for _, kw := range cap.Keywords {
				if kw != "" && strings.Contains(text, kw) {
					score += float64(cap.Priority) * 0.4
				}
			}
		}

		if ctx.Session.LastAgentID != "" && ctx.Session.LastAgentID == a.AgentID {
			score += 10 * 0.3
		}

		maxCap := a.Health.MaxCapacity
		if maxCap <= 0 {
			maxCap = 1
		}
		loadScore := 1.0 - (float64(a.Health.CurrentLoad) / float64(maxCap))
		if loadScore < 0 {
			loadScore = 0
		}
		score += loadScore * 10 * 0.2

		latencyScore := 1.0 / (1.0 + a.Health.AvgLatencyMs/1000.0)
		score += latencyScore * 10 * 0.1

		out = append(out, scored{agent: a, score: score})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].score > out[j].score })
	return out[0].agent, nil
}
