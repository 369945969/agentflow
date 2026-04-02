package sse

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"agentServer/internal/config"
	"agentServer/internal/protocol"
)

// OpencodeMessage represents a message from opencode
// This matches the SSE protocol used by opencode
type OpencodeMessage struct {
	MessageID        string                 `json:"message_id"`
	SessionID        string                 `json:"session_id"`
	Content          string                 `json:"content"`
	IsThinking       bool                   `json:"is_thinking"`
	IsFinal          bool                   `json:"is_final"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	Usage            map[string]int         `json:"usage,omitempty"`
	Status           string                 `json:"status,omitempty"`
	AssistantMessage map[string]interface{} `json:"assistant_message,omitempty"`
}

// SSEClient handles SSE connections to opencode
// This client connects to the external opencode service
type SSEClient struct {
	baseURL     string
	client      *http.Client
	opencodeURL string
}

// NewSSEClient creates a new SSEClient instance
func NewSSEClient(opencodeURL string) *SSEClient {
	return &SSEClient{
		baseURL:     "", // BaseURL not used, opencodeURL is absolute
		opencodeURL: opencodeURL,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Stream connects to opencode via SSE and streams responses
func (c *SSEClient) Stream(ctx context.Context, msg protocol.Message) (<-chan protocol.ServerMessage, <-chan error, error) {
	// Prepare the output channels
	out := make(chan protocol.ServerMessage, 32)
	errCh := make(chan error, 1)

	// Set up SSE URL
	sseURL := fmt.Sprintf("%s/type", c.opencodeURL)

	// Prepare the request payload
	payload := map[string]interface{}{
		"message":           msg.Content.Text,
		"user_id":           msg.UserID,
		"session_id":        msg.SessionID,
		"message_id":        msg.MessageID,
		"thinking_enabled":  msg.ThinkingEnabled,
		"simplified_output": msg.SimplifiedOutput,
		"system_prompt":     msg.SystemPrompt,
		"model_id":          msg.BoundModel,
		"stream":            true,
		"metadata":          msg.Metadata,
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", sseURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for SSE
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request to opencode: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("opencode returned status %d", resp.StatusCode)
	}

	// Initialize SSE stream processing
	go c.processSSEStream(ctx, resp, msg, out, errCh)

	return out, errCh, nil
}

// RequestSummary requests a summary from opencode for simplified output mode
func (c *SSEClient) RequestSummary(ctx context.Context, msg protocol.Message, finalContent string) (<-chan protocol.ServerMessage, <-chan error, error) {
	out := make(chan protocol.ServerMessage, 32)
	errCh := make(chan error, 1)

	// Check if simplified output is enabled
	if !msg.SimplifiedOutput {
		close(out)
		close(errCh)
		return out, errCh, nil
	}

	// Prepare summary request URL
	summaryURL := fmt.Sprintf("%s/summarize", c.opencodeURL)

	// Prepare the request payload
	payload := map[string]interface{}{
		"content":          finalContent,
		"user_id":          msg.UserID,
		"session_id":       msg.SessionID,
		"message_id":       msg.MessageID,
		"original_message": msg.Content.Text,
		"model_id":         msg.BoundModel,
	}

	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal summary payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", summaryURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create summary request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send summary request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("opencode summary returned status %d", resp.StatusCode)
	}

	// Process the summary response
	go c.processSummaryResponse(ctx, resp, msg, out, errCh)

	return out, errCh, nil
}

// processSSEStream handles the SSE response stream from opencode
func (c *SSEClient) processSSEStream(ctx context.Context, resp *http.Response, msg protocol.Message, out chan<- protocol.ServerMessage, errCh chan<- error) {
	defer close(out)
	defer close(errCh)
	defer resp.Body.Close()

	// Create buffered reader for efficient line reading
	reader := bufio.NewReader(resp.Body)

	// Track whether we've sent the final chunk
	finalSent := false
	finalContent := ""

	// Read stream line by line
	for {
		select {
		case <-ctx.Done():
			log.Println("SSE stream context cancelled")
			return
		default:
			// Read a line from the SSE stream
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errCh <- fmt.Errorf("failed to read SSE stream: %w", err)
				}
				return
			}

			// Skip comment lines
			if len(line) < 1 || line[0] == ':' {
				continue
			}

			// Remove the trailing newline
			line = line[:len(line)-1]

			// Handle event lines
			if len(line) > 6 && line[:6] == "event:" {
				// Extract event type
				eventType := line[6:]
				log.Printf("Received SSE event: %s", eventType)
				continue
			}

			// Handle data lines
			if len(line) > 6 && line[:5] == "data:" {
				// Extract data
				data := line[5:]
				if data[0] == ' ' {
					data = data[1:]
				}

				// Unmarshal the JSON data
				var sseMsg OpencodeMessage
				if err := json.Unmarshal([]byte(data), &sseMsg); err != nil {
					log.Printf("Failed to unmarshal SSE data: %v", err)
					continue
				}

				// Skip empty content
				if sseMsg.Content == "" && !sseMsg.IsFinal {
					continue
				}

				// Build protocol.ServerMessage
				response := protocol.ServerMessage{
					Type:      "stream_chunk",
					MessageID: msg.MessageID,
					SessionID: msg.SessionID,
					Payload: protocol.StreamChunkPayload{
						Content:    sseMsg.Content,
						IsThinking: sseMsg.IsThinking,
						IsFinal:    sseMsg.IsFinal,
						SenderID:   "opencode",
						Metadata:   sseMsg.Metadata,
					},
					Timestamp: time.Now().UnixMilli(),
				}

				// Apply user settings filtering
				if !msg.ThinkingEnabled && response.Payload.(protocol.StreamChunkPayload).IsThinking {
					// Skip thinking chunks if thinking is disabled
					continue
				}

				if msg.SimplifiedOutput && !response.Payload.(protocol.StreamChunkPayload).IsFinal {
					// Skip non-final chunks if simplified output is enabled
					continue
				}

				// Send the response
				if !response.Payload.(protocol.StreamChunkPayload).IsThinking || msg.ThinkingEnabled {
					out <- response
				}

				// Track final content
				if response.Payload.(protocol.StreamChunkPayload).IsFinal {
					finalSent = true
					finalContent = sseMsg.Content
				}
			}
		}
	}

	// Request summary if simplified output is enabled and final content was received
	if msg.SimplifiedOutput && finalSent && finalContent != "" {
		summaryOut, summaryErr, err := c.RequestSummary(ctx, msg, finalContent)
		if err == nil {
			// Forward the summary as final content
			for summaryMsg := range summaryOut {
				out <- summaryMsg
			}
			// Forward any summary errors
			if err, ok := <-summaryErr; ok {
				errCh <- fmt.Errorf("summary error: %w", err)
			}
		} else if err != nil {
			errCh <- fmt.Errorf("failed to request summary: %w", err)
		}
	}
}

// processSummaryResponse handles the summary response from opencode
func (c *SSEClient) processSummaryResponse(ctx context.Context, resp *http.Response, msg protocol.Message, out chan<- protocol.ServerMessage, errCh chan<- error) {
	defer close(out)
	defer close(errCh)
	defer resp.Body.Close()

	// Decode the summary response
	var summaryResp struct {
		Summary string `json:"summary"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&summaryResp); err != nil {
		errCh <- fmt.Errorf("failed to decode summary response: %w", err)
		return
	}

	// Create final message
	finalMsg := protocol.ServerMessage{
		Type:      "stream_chunk",
		MessageID: msg.MessageID,
		SessionID: msg.SessionID,
		Payload: protocol.StreamChunkPayload{
			Content:    summaryResp.Summary,
			IsFinal:    true,
			IsThinking: false,
			SenderID:   "opencode",
		},
		Timestamp: time.Now().UnixMilli(),
	}

	// Send the final message
	out <- finalMsg
}

// CreateGlobalSSEClient creates the global SSE client instance
func CreateGlobalSSEClient() (*SSEClient, error) {
	// Load config
	cfg, err := config.Load("")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create client
	client := NewSSEClient(cfg.OpenCode.BaseURL)
	return client, nil
}
