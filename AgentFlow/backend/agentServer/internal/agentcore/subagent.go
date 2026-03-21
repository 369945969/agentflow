package agentcore

type SubAgentConfig struct {
	AgentID         string
	Name            string
	Endpoint        string
	Capabilities    []AgentCapability
	ExecutionConfig AgentExecutionConfig
	Health          HealthStatus
}

type AgentCapability struct {
	SkillID     string
	SkillName   string
	Description string
	Keywords    []string
	InputSchema map[string]interface{}
	Priority    int
}

type AgentExecutionConfig struct {
	ResultMarkers ResultMarkers
	Streaming     StreamingConfig
	TimeoutMs     int
	MaxConcurrent int
}

type ResultMarkers struct {
	FinalPrefix    string
	ThinkingPrefix string
	ErrorPrefix    string
}

type StreamingConfig struct {
	Enabled       bool
	BufferSize    int
	FlushInterval int
}

type HealthStatus struct {
	CurrentLoad     int
	MaxCapacity     int
	AvgLatencyMs    float64
	LastHealthCheck int64
	IsHealthy       bool
}
