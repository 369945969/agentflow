package opencode

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"agentServer/internal/protocol"
)

type Client struct {
	baseURL          string
	backendModelsURL string
	httpClient       *http.Client
}

func NewClient(baseURL string, backendModelsURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		baseURL:          strings.TrimRight(baseURL, "/"),
		backendModelsURL: strings.TrimRight(backendModelsURL, "/"),
		httpClient:       httpClient,
	}
}

func (c *Client) Stream(ctx context.Context, msg protocol.Message, upstreamSessionID string) (string, <-chan protocol.ServerMessage, <-chan error, error) {
	out := make(chan protocol.ServerMessage, 32)
	errCh := make(chan error, 1)

	providers, err := c.fetchProviders(ctx)
	if err != nil {
		close(out)
		close(errCh)
		return "", out, errCh, err
	}

	backendModels, _ := c.fetchBackendModels(ctx)
	selectedBackendModel := pickBackendModel(backendModels, msg.Metadata.ModelID)

	providerID, modelID := resolveProviderModel(providers, selectedBackendModel)
	if providerID == "" || modelID == "" {
		providerID, modelID = fallbackProviderModel(providers)
	}
	if providerID == "" || modelID == "" {
		close(out)
		close(errCh)
		return "", out, errCh, fmt.Errorf("OpenCode provider/model not found")
	}

	sessionID := strings.TrimSpace(upstreamSessionID)
	if sessionID == "" {
		sessionID, err = c.createSession(ctx)
		if err != nil {
			close(out)
			close(errCh)
			return "", out, errCh, err
		}
	}

	go func() {
		defer close(out)
		defer close(errCh)

		eventCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		doneCh := make(chan struct{})
		go func() {
			defer close(doneCh)
			var textBuf strings.Builder
			var reasoningBuf strings.Builder
			c.readEvents(eventCtx, cancel, out, sessionID, msg, &textBuf, &reasoningBuf)
		}()

		duckdbModelID := ""
		if selectedBackendModel != nil {
			duckdbModelID = selectedBackendModel.ID
		} else if msg.Metadata.ModelID != "" {
			duckdbModelID = msg.Metadata.ModelID
		}

		if err := c.sendMessage(ctx, sessionID, providerID, modelID, msg.UserID, duckdbModelID, msg.Content.Text); err != nil {
			errCh <- err
			cancel()
			<-doneCh
			return
		}

		select {
		case <-ctx.Done():
			cancel()
			<-doneCh
			return
		case <-doneCh:
			return
		}
	}()

	return sessionID, out, errCh, nil
}

type backendModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BaseURL   string `json:"base_url"`
	ModelName string `json:"model_name"`
	IsDefault bool   `json:"is_default"`
	Provider  string `json:"provider"`
}

type providerResponse struct {
	All       []providerEntry   `json:"all"`
	Default   map[string]string `json:"default"`
	Connected []string          `json:"connected"`
}

type providerEntry struct {
	ID     string                   `json:"id"`
	Models map[string]providerModel `json:"models"`
}

type providerModel struct {
	API providerModelAPI `json:"api"`
}

type providerModelAPI struct {
	URL string `json:"url"`
}

func (c *Client) fetchProviders(ctx context.Context) (*providerResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/provider", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to OpenCode: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenCode error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var pr providerResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("Invalid JSON format: %w", err)
	}
	return &pr, nil
}

func (c *Client) fetchBackendModels(ctx context.Context) ([]backendModel, error) {
	if c.backendModelsURL == "" {
		return nil, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.backendModelsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend models error (%d)", resp.StatusCode)
	}

	var models []backendModel
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, err
	}
	return models, nil
}

func pickBackendModel(models []backendModel, wantedID string) *backendModel {
	if len(models) == 0 {
		return nil
	}
	if wantedID != "" {
		for i := range models {
			if models[i].ID == wantedID {
				return &models[i]
			}
		}
	}
	for i := range models {
		if models[i].IsDefault {
			return &models[i]
		}
	}
	return &models[0]
}

func resolveProviderModel(pr *providerResponse, bm *backendModel) (string, string) {
	if pr == nil {
		return "", ""
	}
	if bm == nil {
		return "", ""
	}

	providerID := providerForBaseURL(pr, bm.BaseURL)
	if providerID == "" {
		providerID = resolveConnectedProvider(pr)
	}
	if providerID == "" {
		return "", ""
	}

	modelID := bm.ModelName
	if modelID != "" && providerHasModel(pr, providerID, modelID) {
		return providerID, modelID
	}

	modelID = pickModel(pr, providerID)
	return providerID, modelID
}

func fallbackProviderModel(pr *providerResponse) (string, string) {
	if pr == nil {
		return "", ""
	}
	providerID := resolveConnectedProvider(pr)
	if providerID == "" {
		return "", ""
	}
	return providerID, pickModel(pr, providerID)
}

func providerForBaseURL(pr *providerResponse, baseURL string) string {
	if pr == nil || baseURL == "" {
		return ""
	}
	for _, p := range pr.All {
		for _, m := range p.Models {
			if m.API.URL == baseURL {
				return p.ID
			}
		}
	}
	return ""
}

func resolveConnectedProvider(pr *providerResponse) string {
	if pr == nil || len(pr.Connected) == 0 {
		return ""
	}
	return pr.Connected[0]
}

func providerHasModel(pr *providerResponse, providerID string, modelID string) bool {
	if pr == nil || providerID == "" || modelID == "" {
		return false
	}
	for _, p := range pr.All {
		if p.ID != providerID {
			continue
		}
		_, ok := p.Models[modelID]
		return ok
	}
	return false
}

func pickModel(pr *providerResponse, providerID string) string {
	if pr == nil || providerID == "" {
		return ""
	}
	if pr.Default != nil {
		if m := pr.Default[providerID]; m != "" {
			return m
		}
	}
	for _, p := range pr.All {
		if p.ID != providerID {
			continue
		}
		for k := range p.Models {
			return k
		}
	}
	return ""
}

func (c *Client) createSession(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/session", bytes.NewBufferString(`{}`))
	if err != nil {
		return "", fmt.Errorf("Failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to connect to OpenCode: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenCode error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var row struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&row); err != nil {
		return "", fmt.Errorf("Invalid JSON format: %w", err)
	}
	if row.ID == "" {
		return "", fmt.Errorf("OpenCode error (invalid session)")
	}
	return row.ID, nil
}

func (c *Client) sendMessage(ctx context.Context, sessionID, providerID, modelID, userID, duckdbModelID, text string) error {
	body, err := json.Marshal(map[string]interface{}{
		"model": map[string]string{
			"providerID": providerID,
			"modelID":    modelID,
		},
		"variant": userID,
		"parts": []map[string]string{
			{"type": "text", "text": text},
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to create request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/session/"+sessionID+"/message", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}
	if duckdbModelID != "" {
		req.Header.Set("X-Model-ID", duckdbModelID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to connect to OpenCode: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OpenCode error (%d): %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	return nil
}

func (c *Client) readEvents(ctx context.Context, cancel context.CancelFunc, out chan<- protocol.ServerMessage, sessionID string, src protocol.Message, textBuf *strings.Builder, reasoningBuf *strings.Builder) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/event", nil)
	if err != nil {
		out <- protocol.ServerMessage{
			Type:      "error",
			MessageID: src.MessageID,
			SessionID: src.SessionID,
			Timestamp: time.Now().UnixMilli(),
			Payload: protocol.ErrorPayload{
				Code:      "request_error",
				Message:   "Failed to create request",
				Retryable: false,
			},
		}
		cancel()
		return
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		out <- protocol.ServerMessage{
			Type:      "error",
			MessageID: src.MessageID,
			SessionID: src.SessionID,
			Timestamp: time.Now().UnixMilli(),
			Payload: protocol.ErrorPayload{
				Code:      "connect_error",
				Message:   fmt.Sprintf("Failed to connect to OpenCode: %v", err),
				Retryable: true,
			},
		}
		cancel()
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		out <- protocol.ServerMessage{
			Type:      "error",
			MessageID: src.MessageID,
			SessionID: src.SessionID,
			Timestamp: time.Now().UnixMilli(),
			Payload: protocol.ErrorPayload{
				Code:      "opencode_error",
				Message:   fmt.Sprintf("OpenCode error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body))),
				Retryable: true,
			},
		}
		cancel()
		return
	}

	reader := bufio.NewReader(resp.Body)
	assistantMsgID := ""
	gotText := false
	currentEventType := ""
	var dataBuf bytes.Buffer

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			if dataBuf.Len() == 0 {
				currentEventType = ""
				continue
			}

			jsonStr := strings.TrimSpace(dataBuf.String())
			dataBuf.Reset()

			if jsonStr == "" {
				currentEventType = ""
				continue
			}

			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
				currentEventType = ""
				continue
			}

			evType := currentEventType
			currentEventType = ""
			if t, _ := raw["type"].(string); t != "" {
				evType = t
			}

			switch evType {
			case "message.updated":
				sid := getString(raw, "properties", "info", "sessionID")
				role := getString(raw, "properties", "info", "role")
				if sid == sessionID && role == "assistant" {
					assistantMsgID = getString(raw, "properties", "info", "id")
				}
			case "message.part.updated":
				sid := firstNonEmpty(
					getString(raw, "properties", "part", "sessionID"),
					getString(raw, "properties", "sessionID"),
				)
				if sid != sessionID {
					continue
				}
				mid := firstNonEmpty(
					getString(raw, "properties", "part", "messageID"),
					getString(raw, "properties", "messageID"),
				)
				ptype := firstNonEmpty(
					getString(raw, "properties", "part", "type"),
					getString(raw, "properties", "type"),
				)
				text := getString(raw, "properties", "part", "text")
				if text == "" {
					text = getString(raw, "properties", "part", "reasoning")
				}
				delta := firstNonEmpty(
					getStringFlexible(raw, "properties", "delta"),
					getStringFlexible(raw, "properties", "part", "delta"),
					getStringFlexible(raw, "properties", "part", "text_delta"),
					getStringFlexible(raw, "properties", "part", "reasoning_delta"),
				)
				if ptype != "text" && ptype != "reasoning" {
					continue
				}
				if assistantMsgID != "" && mid != assistantMsgID {
					continue
				}
				var targetBuf *strings.Builder
				isThinking := ptype == "reasoning"
				if isThinking {
					targetBuf = reasoningBuf
				} else {
					targetBuf = textBuf
				}
				if targetBuf == nil {
					continue
				}
				if delta != "" {
					targetBuf.WriteString(delta)
				} else {
					prev := targetBuf.String()
					delta = text
					if prev != "" && strings.HasPrefix(text, prev) {
						delta = text[len(prev):]
						targetBuf.Reset()
						targetBuf.WriteString(text)
					} else {
						targetBuf.WriteString(text)
					}
				}
				if delta == "" {
					continue
				}
				log.Printf("[opencode] event message.part.updated session_id=%s source_message_id=%s assistant_message_id=%s part_type=%s delta_len=%d text_len=%d used_native_delta=%t", sessionID, src.MessageID, mid, ptype, len(delta), len(text), getString(raw, "properties", "delta") != "")
				gotText = true
				out <- protocol.ServerMessage{
					Type:      "stream_chunk",
					MessageID: src.MessageID,
					SessionID: src.SessionID,
					Timestamp: time.Now().UnixMilli(),
					Payload: protocol.StreamChunkPayload{
						Content:    delta,
						IsThinking: isThinking,
						IsFinal:    false,
						Metadata: map[string]interface{}{
							"provider_id": getString(raw, "properties", "info", "providerID"),
							"model_id":    getString(raw, "properties", "info", "modelID"),
							"part_type":   ptype,
						},
					},
				}
			case "message.part.delta":
				sid := firstNonEmpty(
					getString(raw, "properties", "part", "sessionID"),
					getString(raw, "properties", "sessionID"),
				)
				if sid != "" && sid != sessionID {
					continue
				}
				mid := firstNonEmpty(
					getString(raw, "properties", "part", "messageID"),
					getString(raw, "properties", "messageID"),
				)
				ptype := firstNonEmpty(
					getString(raw, "properties", "part", "type"),
					getString(raw, "properties", "type"),
				)
				if ptype == "" {
					ptype = "text"
				}
				if ptype != "text" && ptype != "reasoning" {
					continue
				}
				if assistantMsgID != "" && mid != "" && mid != assistantMsgID {
					continue
				}
				delta := firstNonEmpty(
					getStringFlexible(raw, "properties", "delta"),
					getStringFlexible(raw, "properties", "part", "delta"),
					getStringFlexible(raw, "properties", "part", "text_delta"),
					getStringFlexible(raw, "properties", "part", "reasoning_delta"),
					getStringFlexible(raw, "delta"),
				)
				if delta == "" {
					continue
				}
				var targetBuf *strings.Builder
				isThinking := ptype == "reasoning"
				if isThinking {
					targetBuf = reasoningBuf
				} else {
					targetBuf = textBuf
				}
				if targetBuf == nil {
					continue
				}
				log.Printf("[opencode] event message.part.delta session_id=%s source_message_id=%s assistant_message_id=%s part_type=%s delta_len=%d", sessionID, src.MessageID, mid, ptype, len(delta))
				targetBuf.WriteString(delta)
				gotText = true
				out <- protocol.ServerMessage{
					Type:      "stream_chunk",
					MessageID: src.MessageID,
					SessionID: src.SessionID,
					Timestamp: time.Now().UnixMilli(),
					Payload: protocol.StreamChunkPayload{
						Content:    delta,
						IsThinking: isThinking,
						IsFinal:    false,
						Metadata: map[string]interface{}{
							"provider_id": getString(raw, "properties", "info", "providerID"),
							"model_id":    getString(raw, "properties", "info", "modelID"),
							"part_type":   ptype,
						},
					},
				}
			case "session.error":
				sid := getString(raw, "properties", "sessionID")
				if sid != sessionID {
					continue
				}
				msg := getString(raw, "properties", "error", "data", "message")
				if msg == "" {
					msg = getString(raw, "properties", "error", "data")
				}
				if msg == "" {
					msg = "OpenCode session.error"
				}
				if gotText && strings.EqualFold(strings.TrimSpace(msg), "Not Found") {
					log.Printf("[opencode] session.error downgraded to final session_id=%s source_message_id=%s message=%s final_len=%d", sessionID, src.MessageID, msg, textBuf.Len())
					out <- protocol.ServerMessage{
						Type:      "stream_end",
						MessageID: src.MessageID,
						SessionID: src.SessionID,
						Timestamp: time.Now().UnixMilli(),
						Payload: protocol.StreamEndPayload{
							Content: textBuf.String(),
							IsFinal: true,
						},
					}
					cancel()
					return
				}
				out <- protocol.ServerMessage{
					Type:      "error",
					MessageID: src.MessageID,
					SessionID: src.SessionID,
					Timestamp: time.Now().UnixMilli(),
					Payload: protocol.ErrorPayload{
						Code:      "session_error",
						Message:   msg,
						Retryable: true,
					},
				}
				cancel()
				return
			case "session.idle":
				sid := getString(raw, "properties", "sessionID")
				if sid != sessionID {
					continue
				}
				if !gotText {
					continue
				}
				log.Printf("[opencode] event session.idle session_id=%s source_message_id=%s final_len=%d", sessionID, src.MessageID, textBuf.Len())
				out <- protocol.ServerMessage{
					Type:      "stream_end",
					MessageID: src.MessageID,
					SessionID: src.SessionID,
					Timestamp: time.Now().UnixMilli(),
					Payload: protocol.StreamEndPayload{
						Content: textBuf.String(),
						IsFinal: true,
					},
				}
				cancel()
				return
			}
			continue
		}

		if strings.HasPrefix(line, "event:") {
			currentEventType = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}
		if strings.HasPrefix(line, "data:") {
			chunk := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if chunk == "" {
				continue
			}
			if dataBuf.Len() > 0 {
				dataBuf.WriteByte('\n')
			}
			dataBuf.WriteString(chunk)
			continue
		}
	}
}

func getString(m map[string]interface{}, path ...string) string {
	var cur interface{} = m
	for _, p := range path {
		next, ok := cur.(map[string]interface{})
		if !ok {
			return ""
		}
		cur = next[p]
	}
	s, _ := cur.(string)
	return s
}

func getValue(m map[string]interface{}, path ...string) interface{} {
	var cur interface{} = m
	for _, p := range path {
		next, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		cur = next[p]
	}
	return cur
}

func getStringFlexible(m map[string]interface{}, path ...string) string {
	v := getValue(m, path...)
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case map[string]interface{}:
		if s, _ := t["text"].(string); s != "" {
			return s
		}
		if s, _ := t["content"].(string); s != "" {
			return s
		}
		if s, _ := t["delta"].(string); s != "" {
			return s
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
