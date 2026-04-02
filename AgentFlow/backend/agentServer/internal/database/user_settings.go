package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"agentServer/internal/config"
)

// UserSettings represents the settings for a specific user/agent
// These settings are loaded from apiServer (http://76.13.11.164:3000)
type UserSettings struct {
	ThinkingEnabled  bool   `json:"thinking_enabled"`
	SimplifiedOutput bool   `json:"simplified_output"`
	SystemPrompt     string `json:"system_prompt"`
	BoundModel       string `json:"model"`
}

// SettingsLoader handles loading user settings from apiServer
type SettingsLoader struct {
	apiURL     string
	httpClient *http.Client
}

// NewSettingsLoader creates a new SettingsLoader instance
func NewSettingsLoader(apiURL string) *SettingsLoader {
	if apiURL == "" {
		apiURL = "http://localhost:3000"
	}

	return &SettingsLoader{
		apiURL:     apiURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetUserSettings loads user settings from apiServer
// This fetches settings for a specific user/agent by ID
func (s *SettingsLoader) GetUserSettings(ctx context.Context, userID string) (*UserSettings, error) {
	// Construct the API URL to fetch user settings
	url := fmt.Sprintf("%s/api/agents/%s", s.apiURL, userID)

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")

	// Execute the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user settings: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response body
	var settings UserSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate critical settings
	if settings.BoundModel == "" {
		return nil, fmt.Errorf("missing bound model for user %s", userID)
	}

	return &settings, nil
}

// InitializeSettingsLoader creates the global settings loader instance
func InitializeSettingsLoader() (*SettingsLoader, error) {
	// Get the model API URL from config
	cfg, err := config.Load("")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create settings loader
	loader := NewSettingsLoader(cfg.Backend.ModelsAPIURL)
	return loader, nil
}

// Global settings loader instance
var settingsLoader *SettingsLoader

// InitDatabase initializes the database integration
func InitDatabase() error {
	var err error
	settingsLoader, err = InitializeSettingsLoader()
	if err != nil {
		return fmt.Errorf("failed to initialize settings loader: %w", err)
	}

	log.Println("✅ Database integration initialized")
	return nil
}

// GetUserSettings loads user settings from the global settings loader
// This is the public interface to fetch user-specific settings
func GetUserSettings(ctx context.Context, userID string) (*UserSettings, error) {
	if settingsLoader == nil {
		return nil, fmt.Errorf("settings loader not initialized")
	}

	return settingsLoader.GetUserSettings(ctx, userID)
}
