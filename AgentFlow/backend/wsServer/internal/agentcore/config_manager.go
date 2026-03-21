package agentcore

import (
	"errors"
	"sync"
	"time"
)

type UserProfile struct {
	UserID           string
	CurrentSessionID string
	Config           UserInteractionConfig
	CreatedAt        int64
	UpdatedAt        int64
}

type UserInteractionConfig struct {
	StreamMode     string
	EnableThinking bool
	ReturnStrategy ReturnStrategy
	TimeoutMs      int
	MaxRetries     int
}

type ReturnStrategy struct {
	Type              string
	FinalResultMarker string
	ChunkSize         int
	DebounceMs        int
}

type ConfigManager struct {
	mu         sync.RWMutex
	byUser     map[string]*UserProfile
	defaultCfg UserInteractionConfig
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		byUser: make(map[string]*UserProfile),
		defaultCfg: UserInteractionConfig{
			StreamMode:     "realtime",
			EnableThinking: true,
			TimeoutMs:      30000,
			MaxRetries:     1,
			ReturnStrategy: ReturnStrategy{
				Type:              "immediate",
				FinalResultMarker: "[FINAL]",
				ChunkSize:         1024,
				DebounceMs:        50,
			},
		},
	}
}

func (cm *ConfigManager) SetDefault(cfg UserInteractionConfig) {
	cm.mu.Lock()
	cm.defaultCfg = cfg
	cm.mu.Unlock()
}

func (cm *ConfigManager) DefaultConfig() UserInteractionConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.defaultCfg
}

func (cm *ConfigManager) GetUserProfile(userID string) (*UserProfile, error) {
	cm.mu.RLock()
	p, ok := cm.byUser[userID]
	cm.mu.RUnlock()
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (cm *ConfigManager) GetDefaultProfile(userID string) *UserProfile {
	now := time.Now().Unix()
	return &UserProfile{
		UserID:    userID,
		Config:    cm.DefaultConfig(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (cm *ConfigManager) UpdateUserProfile(profile *UserProfile) error {
	if profile == nil || profile.UserID == "" {
		return errors.New("invalid profile")
	}
	now := time.Now().Unix()
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if profile.CreatedAt == 0 {
		profile.CreatedAt = now
	}
	profile.UpdatedAt = now
	cm.byUser[profile.UserID] = profile
	return nil
}

func (cm *ConfigManager) GetOrCreate(userID string) *UserProfile {
	now := time.Now().Unix()

	cm.mu.RLock()
	p, ok := cm.byUser[userID]
	cm.mu.RUnlock()
	if ok {
		cm.mu.Lock()
		p.UpdatedAt = now
		cm.mu.Unlock()
		return p
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()
	if p, ok := cm.byUser[userID]; ok {
		p.UpdatedAt = now
		return p
	}
	p = &UserProfile{
		UserID:    userID,
		Config:    cm.defaultCfg,
		CreatedAt: now,
		UpdatedAt: now,
	}
	cm.byUser[userID] = p
	return p
}
