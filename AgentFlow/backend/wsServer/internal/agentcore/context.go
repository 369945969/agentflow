package agentcore

import "wsServer/internal/protocol"

type ExecutionContext struct {
	Message         protocol.Message
	Session         Session
	UserProfile     UserProfile
	ConnectionID    string
	TargetAgent     SubAgentConfig
	RoleIdentity    *RoleIdentity
	SystemPrompt    string
	EffectiveConfig UserInteractionConfig
}

type RoleIdentity struct {
	UserID          string
	RoleDescription string
	Skills          []Skill
	Persona         string
}

type GroupInfo struct {
	GroupID   string
	GroupName string
	Members   []GroupMember
	CreatorID string
}

type GroupMember struct {
	UserID      string
	Role        string
	Skills      []Skill
	Description string
	Priority    int
	IsActive    bool
	Persona     string
}

type Skill struct {
	SkillID     string
	SkillName   string
	Keywords    []string
	Description string
	Priority    int
}

func mergeUserConfig(defaultCfg UserInteractionConfig, userCfg UserInteractionConfig, override *protocol.UserInteractionConfig) UserInteractionConfig {
	out := defaultCfg

	if userCfg.StreamMode != "" {
		out.StreamMode = userCfg.StreamMode
	}
	out.EnableThinking = userCfg.EnableThinking
	if userCfg.TimeoutMs != 0 {
		out.TimeoutMs = userCfg.TimeoutMs
	}
	if userCfg.MaxRetries != 0 {
		out.MaxRetries = userCfg.MaxRetries
	}
	if userCfg.ReturnStrategy.Type != "" {
		out.ReturnStrategy.Type = userCfg.ReturnStrategy.Type
	}
	if userCfg.ReturnStrategy.FinalResultMarker != "" {
		out.ReturnStrategy.FinalResultMarker = userCfg.ReturnStrategy.FinalResultMarker
	}
	if userCfg.ReturnStrategy.ChunkSize != 0 {
		out.ReturnStrategy.ChunkSize = userCfg.ReturnStrategy.ChunkSize
	}
	if userCfg.ReturnStrategy.DebounceMs != 0 {
		out.ReturnStrategy.DebounceMs = userCfg.ReturnStrategy.DebounceMs
	}

	if override == nil {
		return out
	}

	if override.StreamMode != "" {
		out.StreamMode = override.StreamMode
	}
	out.EnableThinking = override.EnableThinking
	if override.TimeoutMs != 0 {
		out.TimeoutMs = override.TimeoutMs
	}
	if override.MaxRetries != 0 {
		out.MaxRetries = override.MaxRetries
	}
	if override.ReturnStrategy.Type != "" {
		out.ReturnStrategy.Type = override.ReturnStrategy.Type
	}
	if override.ReturnStrategy.FinalResultMarker != "" {
		out.ReturnStrategy.FinalResultMarker = override.ReturnStrategy.FinalResultMarker
	}
	if override.ReturnStrategy.ChunkSize != 0 {
		out.ReturnStrategy.ChunkSize = override.ReturnStrategy.ChunkSize
	}
	if override.ReturnStrategy.DebounceMs != 0 {
		out.ReturnStrategy.DebounceMs = override.ReturnStrategy.DebounceMs
	}

	return out
}
