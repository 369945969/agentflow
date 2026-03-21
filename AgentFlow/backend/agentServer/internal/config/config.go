package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Server ServerConfig `json:"server"`

	Backend  BackendConfig  `json:"backend"`
	OpenCode OpenCodeConfig `json:"opencode"`

	Agents []AgentConfig `json:"agents"`

	DefaultUserConfig UserInteractionConfig `json:"default_user_config"`
	Log               LogConfig             `json:"log"`
	WebSocket         WebSocketConfig       `json:"websocket"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type BackendConfig struct {
	ModelsAPIURL string `json:"models_api_url"`
}

type OpenCodeConfig struct {
	BaseURL string `json:"base_url"`
}

type AgentConfig struct {
	AgentID         string               `json:"agent_id"`
	Name            string               `json:"name"`
	Endpoint        string               `json:"endpoint"`
	Capabilities    []AgentCapability    `json:"capabilities"`
	ExecutionConfig AgentExecutionConfig `json:"execution_config"`
	Health          HealthStatus         `json:"health"`
}

type AgentCapability struct {
	SkillID     string                 `json:"skill_id"`
	SkillName   string                 `json:"skill_name"`
	Description string                 `json:"description"`
	Keywords    []string               `json:"keywords"`
	InputSchema map[string]interface{} `json:"input_schema"`
	Priority    int                    `json:"priority"`
}

type AgentExecutionConfig struct {
	ResultMarkers ResultMarkers   `json:"result_markers"`
	Streaming     StreamingConfig `json:"streaming"`
	TimeoutMs     int             `json:"timeout_ms"`
	MaxConcurrent int             `json:"max_concurrent_reqs"`
}

type ResultMarkers struct {
	FinalPrefix    string `json:"final_prefix"`
	ThinkingPrefix string `json:"thinking_prefix"`
	ErrorPrefix    string `json:"error_prefix"`
}

type StreamingConfig struct {
	Enabled       bool `json:"enabled"`
	BufferSize    int  `json:"buffer_size"`
	FlushInterval int  `json:"flush_interval"`
}

type HealthStatus struct {
	CurrentLoad     int     `json:"current_load"`
	MaxCapacity     int     `json:"max_capacity"`
	AvgLatencyMs    float64 `json:"avg_latency_ms"`
	LastHealthCheck int64   `json:"last_health_check"`
	IsHealthy       bool    `json:"is_healthy"`
}

type LogConfig struct {
	Dir string `json:"dir"`
}

type WebSocketConfig struct {
	PingIntervalSec int `json:"ping_interval_sec"`
	PongWaitSec     int `json:"pong_wait_sec"`
	SendQueueSize   int `json:"send_queue_size"`
}

type UserInteractionConfig struct {
	StreamMode     string         `json:"stream_mode"`
	EnableThinking bool           `json:"enable_thinking"`
	ReturnStrategy ReturnStrategy `json:"return_strategy"`
	TimeoutMs      int            `json:"timeout_ms"`
	MaxRetries     int            `json:"max_retries"`
}

type ReturnStrategy struct {
	Type              string `json:"type"`
	FinalResultMarker string `json:"final_result_marker"`
	ChunkSize         int    `json:"chunk_size"`
	DebounceMs        int    `json:"debounce_ms"`
}

func DefaultPath() string {
	if p := strings.TrimSpace(os.Getenv("WS_SERVER_CONFIG")); p != "" {
		return p
	}
	return filepath.Join("config", "config.json")
}

func Load(path string) (*Config, error) {
	if strings.TrimSpace(path) == "" {
		path = DefaultPath()
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 3002
	}
	if cfg.OpenCode.BaseURL == "" {
		cfg.OpenCode.BaseURL = "http://localhost:3001"
	}
	if cfg.Backend.ModelsAPIURL == "" {
		cfg.Backend.ModelsAPIURL = "http://localhost:3000/api/models/"
	}
	if cfg.Log.Dir == "" {
		cfg.Log.Dir = "logs"
	}
	if cfg.WebSocket.PingIntervalSec == 0 {
		cfg.WebSocket.PingIntervalSec = 30
	}
	if cfg.WebSocket.PongWaitSec == 0 {
		cfg.WebSocket.PongWaitSec = 60
	}
	if cfg.WebSocket.SendQueueSize == 0 {
		cfg.WebSocket.SendQueueSize = 256
	}
	if len(cfg.Agents) == 0 {
		cfg.Agents = []AgentConfig{
			{AgentID: "opencode-1", Name: "OpenCode", Endpoint: cfg.OpenCode.BaseURL},
		}
	}

	return &cfg, nil
}
