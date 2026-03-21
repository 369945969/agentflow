package protocol

type Message struct {
	MessageID string          `json:"message_id"`
	SessionID string          `json:"session_id"`
	GroupID   string          `json:"group_id,omitempty"`
	UserID    string          `json:"user_id"`
	Type      string          `json:"type"`
	Content   MessageContent  `json:"content"`
	Metadata  MessageMetadata `json:"metadata"`
	Timestamp int64           `json:"timestamp"`
}

type MessageContent struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Mentions    []string     `json:"mentions,omitempty"`
}

type Attachment struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}

type MessageMetadata struct {
	OverrideConfig *UserInteractionConfig `json:"override_config,omitempty"`
	RequiredSkills []string               `json:"required_skills,omitempty"`
	TargetAgentID  string                 `json:"target_agent_id,omitempty"`
	ModelID        string                 `json:"model_id,omitempty"`
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

type ServerMessage struct {
	Type      string      `json:"type"`
	MessageID string      `json:"message_id,omitempty"`
	SessionID string      `json:"session_id,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

type StreamChunkPayload struct {
	Content    string                 `json:"content"`
	IsThinking bool                   `json:"is_thinking"`
	IsFinal    bool                   `json:"is_final"`
	Progress   int                    `json:"progress,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type StreamEndPayload struct {
	Content string         `json:"content"`
	IsFinal bool           `json:"is_final"`
	Usage   map[string]int `json:"usage,omitempty"`
}

type ErrorPayload struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}
