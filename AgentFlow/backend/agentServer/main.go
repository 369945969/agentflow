package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"agentServer/internal/agentcore"
	"agentServer/internal/config"
	"agentServer/internal/protocol"
	"agentServer/internal/ws"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatal(err)
	}

	agents := make([]agentcore.SubAgentConfig, 0, len(cfg.Agents))
	for _, a := range cfg.Agents {
		caps := make([]agentcore.AgentCapability, 0, len(a.Capabilities))
		for _, c := range a.Capabilities {
			caps = append(caps, agentcore.AgentCapability{
				SkillID:     c.SkillID,
				SkillName:   c.SkillName,
				Description: c.Description,
				Keywords:    c.Keywords,
				InputSchema: c.InputSchema,
				Priority:    c.Priority,
			})
		}
		agents = append(agents, agentcore.SubAgentConfig{
			AgentID:      a.AgentID,
			Name:         a.Name,
			Endpoint:     a.Endpoint,
			Capabilities: caps,
			ExecutionConfig: agentcore.AgentExecutionConfig{
				ResultMarkers: agentcore.ResultMarkers{
					FinalPrefix:    a.ExecutionConfig.ResultMarkers.FinalPrefix,
					ThinkingPrefix: a.ExecutionConfig.ResultMarkers.ThinkingPrefix,
					ErrorPrefix:    a.ExecutionConfig.ResultMarkers.ErrorPrefix,
				},
				Streaming: agentcore.StreamingConfig{
					Enabled:       a.ExecutionConfig.Streaming.Enabled,
					BufferSize:    a.ExecutionConfig.Streaming.BufferSize,
					FlushInterval: a.ExecutionConfig.Streaming.FlushInterval,
				},
				TimeoutMs:     a.ExecutionConfig.TimeoutMs,
				MaxConcurrent: a.ExecutionConfig.MaxConcurrent,
			},
			Health: agentcore.HealthStatus{
				CurrentLoad:     a.Health.CurrentLoad,
				MaxCapacity:     a.Health.MaxCapacity,
				AvgLatencyMs:    a.Health.AvgLatencyMs,
				LastHealthCheck: a.Health.LastHealthCheck,
				IsHealthy:       a.Health.IsHealthy,
			},
		})
	}

	core := agentcore.NewAgentCore(agents, cfg.Backend.ModelsAPIURL, nil)
	core.SetDefaultUserConfig(agentcore.UserInteractionConfig{
		StreamMode:     cfg.DefaultUserConfig.StreamMode,
		EnableThinking: cfg.DefaultUserConfig.EnableThinking,
		TimeoutMs:      cfg.DefaultUserConfig.TimeoutMs,
		MaxRetries:     cfg.DefaultUserConfig.MaxRetries,
		ReturnStrategy: agentcore.ReturnStrategy{
			Type:              cfg.DefaultUserConfig.ReturnStrategy.Type,
			FinalResultMarker: cfg.DefaultUserConfig.ReturnStrategy.FinalResultMarker,
			ChunkSize:         cfg.DefaultUserConfig.ReturnStrategy.ChunkSize,
			DebounceMs:        cfg.DefaultUserConfig.ReturnStrategy.DebounceMs,
		},
	})
	core.SetLogBasePath(cfg.Log.Dir)

	gateway := ws.NewGateway(upgrader, core, func(ctx context.Context, connID string, msg protocol.Message) {
		core.HandleMessage(ctx, connID, msg)
	},
		ws.WithHeartbeat(time.Duration(cfg.WebSocket.PingIntervalSec)*time.Second, time.Duration(cfg.WebSocket.PongWaitSec)*time.Second),
		ws.WithSendQueueSize(cfg.WebSocket.SendQueueSize),
	)
	core.SetEmitter(gateway)
	http.Handle("/ws", gateway)
	port := strconv.Itoa(cfg.Server.Port)
	fmt.Printf("WebSocket Server starting on ws://localhost:%s/ws\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
