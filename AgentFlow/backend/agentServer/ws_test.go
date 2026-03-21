package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"agentServer/internal/agentcore"
	"agentServer/internal/protocol"
	wsInternal "agentServer/internal/ws"
)

func TestWebSocketForwarding(t *testing.T) {
	backendMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/models/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"duckdb-1","name":"default","base_url":"http://mock-api","model_name":"mock-model","provider":"mock","is_default":true}]`))
	}))
	defer backendMock.Close()

	events := make(chan string, 16)
	defer close(events)
	var sessionMu sync.Mutex
	sessionSeq := 0

	openCodeMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/provider":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"all": []map[string]interface{}{
					{
						"id": "mock-provider",
						"models": map[string]interface{}{
							"mock-model": map[string]interface{}{
								"api": map[string]interface{}{
									"url": "http://mock-api",
								},
							},
						},
					},
				},
				"default": map[string]string{
					"mock-provider": "mock-model",
				},
				"connected": []string{"mock-provider"},
			})
			return
		case r.Method == http.MethodPost && r.URL.Path == "/session":
			sessionMu.Lock()
			sessionSeq++
			sid := "session-" + strconv.Itoa(sessionSeq)
			sessionMu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"` + sid + `"}`))
			return
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/session/") && strings.HasSuffix(r.URL.Path, "/message"):
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			if len(parts) != 3 {
				http.NotFound(w, r)
				return
			}
			sid := parts[1]
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))

			events <- `{"type":"message.updated","properties":{"info":{"sessionID":"` + sid + `","role":"assistant","id":"assist-1","providerID":"mock-provider","modelID":"mock-model"}}}`
			events <- `{"type":"message.part.updated","properties":{"part":{"sessionID":"` + sid + `","messageID":"assist-1","type":"text","text":"hello"}}}`
			events <- `{"type":"session.idle","properties":{"sessionID":"` + sid + `"}}`
			return
		case r.Method == http.MethodGet && r.URL.Path == "/event":
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "streaming not supported", http.StatusInternalServerError)
				return
			}
			for {
				select {
				case <-r.Context().Done():
					return
				case ev, ok := <-events:
					if !ok {
						return
					}
					_, _ = w.Write([]byte("data: " + ev + "\n\n"))
					flusher.Flush()
				}
			}
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer openCodeMock.Close()

	core := agentcore.NewAgentCore(
		[]agentcore.SubAgentConfig{{AgentID: "opencode-1", Endpoint: openCodeMock.URL}},
		backendMock.URL+"/api/models/",
		nil,
	)
	core.SetLogBasePath(t.TempDir())
	gateway := wsInternal.NewGateway(upgrader, func(ctx context.Context, connID string, msg protocol.Message) {
		core.HandleMessage(ctx, connID, msg)
	})
	core.SetEmitter(gateway)

	testServer := httptest.NewServer(gateway)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	u.Scheme = "ws"
	u.Path = "/ws"
	t.Logf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer c.Close()

	testMsg := protocol.Message{
		MessageID: "msg-1",
		SessionID: "",
		UserID:    "test_user",
		Type:      "text",
		Content: protocol.MessageContent{
			Text: "你好",
		},
	}
	payload, _ := json.Marshal(testMsg)
	if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	receivedChunk := false
	receivedEnd := false

	c.SetReadDeadline(time.Now().Add(10 * time.Second))

	for i := 0; i < 20; i++ {
		_, message, err := c.ReadMessage()
		if err != nil {
			t.Logf("Read loop ended: %v", err)
			break
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(message, &resp); err != nil {
			t.Errorf("Invalid JSON received: %s", string(message))
			continue
		}

		msgType, _ := resp["type"].(string)
		t.Logf("Received message type: %s", msgType)

		switch msgType {
		case "stream_chunk":
			receivedChunk = true
		case "stream_end":
			receivedEnd = true
			goto EndLoop
		case "error":
			t.Fatalf("Server returned error: %v", resp["payload"])
		}
	}

EndLoop:
	if !receivedChunk && !receivedEnd {
		t.Error("Did not receive any stream_chunk or stream_end messages from server")
	}

	if receivedChunk {
		t.Log("Successfully received stream chunk")
	}
	if receivedEnd {
		t.Log("Successfully received stream end")
	}
}

func TestWebSocketValidationRequiresUserOrGroup(t *testing.T) {
	gateway := wsInternal.NewGateway(upgrader, func(ctx context.Context, connID string, msg protocol.Message) {})
	testServer := httptest.NewServer(gateway)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	u.Scheme = "ws"
	u.Path = "/ws"

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer c.Close()

	testMsg := protocol.Message{
		MessageID: "msg-1",
		SessionID: "",
		UserID:    "",
		Type:      "text",
		Content: protocol.MessageContent{
			Text: "你好",
		},
	}
	payload, _ := json.Marshal(testMsg)
	if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, message, err := c.ReadMessage()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(message, &resp); err != nil {
		t.Fatalf("Invalid JSON received: %s", string(message))
	}

	if resp["type"] != "error" {
		t.Fatalf("expected error message, got: %v", resp["type"])
	}
	payloadMap, _ := resp["payload"].(map[string]interface{})
	if payloadMap == nil {
		t.Fatalf("expected error payload, got: %v", resp["payload"])
	}
	if payloadMap["code"] != "validation_error" {
		t.Fatalf("expected validation_error, got: %v", payloadMap["code"])
	}
}

func TestWebSocketGroupAllowsMultipleUsersSameGroup(t *testing.T) {
	backendMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/models/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"duckdb-1","name":"default","base_url":"http://mock-api","model_name":"mock-model","provider":"mock","is_default":true}]`))
	}))
	defer backendMock.Close()

	events := make(chan string, 16)
	defer close(events)
	var sessionMu sync.Mutex
	sessionSeq := 0

	openCodeMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/provider":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"all": []map[string]interface{}{
					{
						"id": "mock-provider",
						"models": map[string]interface{}{
							"mock-model": map[string]interface{}{
								"api": map[string]interface{}{
									"url": "http://mock-api",
								},
							},
						},
					},
				},
				"default": map[string]string{
					"mock-provider": "mock-model",
				},
				"connected": []string{"mock-provider"},
			})
			return
		case r.Method == http.MethodPost && r.URL.Path == "/session":
			sessionMu.Lock()
			sessionSeq++
			sid := "session-" + strconv.Itoa(sessionSeq)
			sessionMu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"` + sid + `"}`))
			return
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/session/") && strings.HasSuffix(r.URL.Path, "/message"):
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			if len(parts) != 3 {
				http.NotFound(w, r)
				return
			}
			sid := parts[1]
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))

			events <- `{"type":"message.updated","properties":{"info":{"sessionID":"` + sid + `","role":"assistant","id":"assist-1","providerID":"mock-provider","modelID":"mock-model"}}}`
			events <- `{"type":"message.part.updated","properties":{"part":{"sessionID":"` + sid + `","messageID":"assist-1","type":"text","text":"hello"}}}`
			events <- `{"type":"session.idle","properties":{"sessionID":"` + sid + `"}}`
			return
		case r.Method == http.MethodGet && r.URL.Path == "/event":
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "streaming not supported", http.StatusInternalServerError)
				return
			}
			for {
				select {
				case <-r.Context().Done():
					return
				case ev, ok := <-events:
					if !ok {
						return
					}
					_, _ = w.Write([]byte("data: " + ev + "\n\n"))
					flusher.Flush()
				}
			}
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer openCodeMock.Close()

	core := agentcore.NewAgentCore(
		[]agentcore.SubAgentConfig{{AgentID: "opencode-1", Endpoint: openCodeMock.URL}},
		backendMock.URL+"/api/models/",
		nil,
	)
	core.SetLogBasePath(t.TempDir())
	gateway := wsInternal.NewGateway(upgrader, func(ctx context.Context, connID string, msg protocol.Message) {
		core.HandleMessage(ctx, connID, msg)
	})
	core.SetEmitter(gateway)

	testServer := httptest.NewServer(gateway)
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	u.Scheme = "ws"
	u.Path = "/ws"

	c1, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Dial c1 error: %v", err)
	}
	defer c1.Close()

	c2, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Dial c2 error: %v", err)
	}
	defer c2.Close()

	readUntilEndOrError := func(c *websocket.Conn) (bool, string) {
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		for i := 0; i < 50; i++ {
			_, message, err := c.ReadMessage()
			if err != nil {
				return false, err.Error()
			}
			var resp map[string]interface{}
			if err := json.Unmarshal(message, &resp); err != nil {
				continue
			}
			msgType, _ := resp["type"].(string)
			if msgType == "error" {
				payloadMap, _ := resp["payload"].(map[string]interface{})
				if payloadMap == nil {
					return false, "error"
				}
				code, _ := payloadMap["code"].(string)
				if code == "" {
					code = "error"
				}
				return false, code
			}
			if msgType == "stream_end" {
				return true, ""
			}
		}
		return false, "timeout"
	}

	groupID := "group-1"

	msg1 := protocol.Message{
		MessageID: "msg-1",
		SessionID: "",
		GroupID:   groupID,
		UserID:    "user_a",
		Type:      "text",
		Content: protocol.MessageContent{
			Text: "你好",
		},
	}
	payload1, _ := json.Marshal(msg1)
	if err := c1.WriteMessage(websocket.TextMessage, payload1); err != nil {
		t.Fatalf("Write c1 error: %v", err)
	}

	if ok, code := readUntilEndOrError(c1); !ok {
		t.Fatalf("c1 expected stream_end, got error: %s", code)
	}

	msg2 := protocol.Message{
		MessageID: "msg-2",
		SessionID: "",
		GroupID:   groupID,
		UserID:    "user_b",
		Type:      "text",
		Content: protocol.MessageContent{
			Text: "再来一次",
		},
	}
	payload2, _ := json.Marshal(msg2)
	if err := c2.WriteMessage(websocket.TextMessage, payload2); err != nil {
		t.Fatalf("Write c2 error: %v", err)
	}

	if ok, code := readUntilEndOrError(c2); !ok {
		t.Fatalf("c2 expected stream_end, got error: %s", code)
	}
	if ok, code := readUntilEndOrError(c1); !ok {
		t.Fatalf("c1 expected stream_end, got error: %s", code)
	}
}

func TestSingleChatThinkingDisabledDoesNotForwardThink(t *testing.T) {
	backendMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/models/":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"duckdb-1","name":"default","base_url":"http://mock-api","model_name":"mock-model","provider":"mock","is_default":true}]`))
			return
		case r.Method == http.MethodGet && r.URL.Path == "/api/agents/":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"user_a","model":"duckdb-1","system_prompt":"SP","thinking_enabled":false,"simplified_output":false}]`))
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer backendMock.Close()

	events := make(chan string, 16)
	defer close(events)
	var sessionMu sync.Mutex
	sessionSeq := 0

	openCodeMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/provider":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"all": []map[string]interface{}{
					{
						"id": "mock-provider",
						"models": map[string]interface{}{
							"mock-model": map[string]interface{}{
								"api": map[string]interface{}{
									"url": "http://mock-api",
								},
							},
						},
					},
				},
				"default": map[string]string{
					"mock-provider": "mock-model",
				},
				"connected": []string{"mock-provider"},
			})
			return
		case r.Method == http.MethodPost && r.URL.Path == "/session":
			sessionMu.Lock()
			sessionSeq++
			sid := "session-" + strconv.Itoa(sessionSeq)
			sessionMu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"` + sid + `"}`))
			return
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/session/") && strings.HasSuffix(r.URL.Path, "/message"):
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			if len(parts) != 3 {
				http.NotFound(w, r)
				return
			}
			sid := parts[1]
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))

			events <- `{"type":"message.updated","properties":{"info":{"sessionID":"` + sid + `","role":"assistant","id":"assist-1","providerID":"mock-provider","modelID":"mock-model"}}}`
			events <- `{"type":"message.part.updated","properties":{"part":{"sessionID":"` + sid + `","messageID":"assist-1","type":"text","text":"<think>abc</think>hello"}}}`
			events <- `{"type":"session.idle","properties":{"sessionID":"` + sid + `"}}`
			return
		case r.Method == http.MethodGet && r.URL.Path == "/event":
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "streaming not supported", http.StatusInternalServerError)
				return
			}
			for {
				select {
				case <-r.Context().Done():
					return
				case ev, ok := <-events:
					if !ok {
						return
					}
					_, _ = w.Write([]byte("data: " + ev + "\n\n"))
					flusher.Flush()
				}
			}
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer openCodeMock.Close()

	core := agentcore.NewAgentCore(
		[]agentcore.SubAgentConfig{{AgentID: "opencode-1", Endpoint: openCodeMock.URL}},
		backendMock.URL+"/api/models/",
		nil,
	)
	logDir := t.TempDir()
	core.SetLogBasePath(logDir)
	gateway := wsInternal.NewGateway(upgrader, func(ctx context.Context, connID string, msg protocol.Message) {
		core.HandleMessage(ctx, connID, msg)
	})
	core.SetEmitter(gateway)

	testServer := httptest.NewServer(gateway)
	defer testServer.Close()

	u, _ := url.Parse(testServer.URL)
	u.Scheme = "ws"
	u.Path = "/ws"
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer c.Close()

	testMsg := protocol.Message{
		MessageID: "msg-1",
		SessionID: "sess-1",
		UserID:    "user_a",
		Type:      "text",
		Content: protocol.MessageContent{
			Text: "你好",
		},
	}
	payload, _ := json.Marshal(testMsg)
	if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	var sawChunk bool
	var sawThink bool
	var sawEnd bool

	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	for i := 0; i < 10; i++ {
		_, message, err := c.ReadMessage()
		if err != nil {
			t.Fatalf("Read error: %v", err)
		}
		var resp map[string]interface{}
		_ = json.Unmarshal(message, &resp)
		if resp["type"] == "stream_chunk" {
			sawChunk = true
			pl, _ := resp["payload"].(map[string]interface{})
			if pl != nil {
				if b, ok := pl["is_thinking"].(bool); ok && b {
					sawThink = true
				}
				if s, ok := pl["content"].(string); ok && strings.Contains(s, "<think>") {
					sawThink = true
				}
			}
		}
		if resp["type"] == "stream_end" {
			sawEnd = true
			pl, _ := resp["payload"].(map[string]interface{})
			if pl != nil {
				if s, ok := pl["content"].(string); ok && s != "hello" {
					t.Fatalf("expected final content hello, got: %s", s)
				}
			}
			break
		}
		if resp["type"] == "error" {
			t.Fatalf("Server returned error: %v", resp["payload"])
		}
	}

	if !sawChunk || !sawEnd {
		t.Fatalf("expected stream_chunk and stream_end, got chunk=%v end=%v", sawChunk, sawEnd)
	}
	if sawThink {
		t.Fatalf("expected no thinking forwarded")
	}

	logPath := filepath.Join(logDir, "user_a", "sess-1.log")
	raw, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log error: %v", err)
	}
	if !strings.Contains(string(raw), "[THINKING] abc") {
		t.Fatalf("expected thinking saved in log")
	}
	if !strings.Contains(string(raw), "[CHUNK] hello") {
		t.Fatalf("expected chunk saved in log")
	}
}

func TestSingleChatSimplifiedOutputOnlySendsSummary(t *testing.T) {
	backendMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/models/":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"duckdb-1","name":"default","base_url":"http://mock-api","model_name":"mock-model","provider":"mock","is_default":true}]`))
			return
		case r.Method == http.MethodGet && r.URL.Path == "/api/agents/":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"user_a","model":"duckdb-1","system_prompt":"SP","thinking_enabled":true,"simplified_output":true}]`))
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer backendMock.Close()

	events := make(chan string, 32)
	defer close(events)
	var sessionMu sync.Mutex
	sessionSeq := 0

	openCodeMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/provider":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"all": []map[string]interface{}{
					{
						"id": "mock-provider",
						"models": map[string]interface{}{
							"mock-model": map[string]interface{}{
								"api": map[string]interface{}{
									"url": "http://mock-api",
								},
							},
						},
					},
				},
				"default": map[string]string{
					"mock-provider": "mock-model",
				},
				"connected": []string{"mock-provider"},
			})
			return
		case r.Method == http.MethodPost && r.URL.Path == "/session":
			sessionMu.Lock()
			sessionSeq++
			sid := "session-" + strconv.Itoa(sessionSeq)
			sessionMu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"` + sid + `"}`))
			return
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/session/") && strings.HasSuffix(r.URL.Path, "/message"):
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			if len(parts) != 3 {
				http.NotFound(w, r)
				return
			}
			sid := parts[1]
			body, _ := io.ReadAll(r.Body)
			isSummary := strings.Contains(string(body), "[完整回答]")

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))

			events <- `{"type":"message.updated","properties":{"info":{"sessionID":"` + sid + `","role":"assistant","id":"assist-1","providerID":"mock-provider","modelID":"mock-model"}}}`
			if isSummary {
				events <- `{"type":"message.part.updated","properties":{"part":{"sessionID":"` + sid + `","messageID":"assist-1","type":"text","text":"summary"}}}`
			} else {
				events <- `{"type":"message.part.updated","properties":{"part":{"sessionID":"` + sid + `","messageID":"assist-1","type":"text","text":"<think>abc</think>raw"}}}`
			}
			events <- `{"type":"session.idle","properties":{"sessionID":"` + sid + `"}}`
			return
		case r.Method == http.MethodGet && r.URL.Path == "/event":
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "streaming not supported", http.StatusInternalServerError)
				return
			}
			for {
				select {
				case <-r.Context().Done():
					return
				case ev, ok := <-events:
					if !ok {
						return
					}
					_, _ = w.Write([]byte("data: " + ev + "\n\n"))
					flusher.Flush()
				}
			}
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer openCodeMock.Close()

	core := agentcore.NewAgentCore(
		[]agentcore.SubAgentConfig{{AgentID: "opencode-1", Endpoint: openCodeMock.URL}},
		backendMock.URL+"/api/models/",
		nil,
	)
	logDir := t.TempDir()
	core.SetLogBasePath(logDir)
	gateway := wsInternal.NewGateway(upgrader, func(ctx context.Context, connID string, msg protocol.Message) {
		core.HandleMessage(ctx, connID, msg)
	})
	core.SetEmitter(gateway)

	testServer := httptest.NewServer(gateway)
	defer testServer.Close()

	u, _ := url.Parse(testServer.URL)
	u.Scheme = "ws"
	u.Path = "/ws"
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Dial error: %v", err)
	}
	defer c.Close()

	testMsg := protocol.Message{
		MessageID: "msg-1",
		SessionID: "sess-1",
		UserID:    "user_a",
		Type:      "text",
		Content: protocol.MessageContent{
			Text: "你好",
		},
	}
	payload, _ := json.Marshal(testMsg)
	if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, message, err := c.ReadMessage()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(message, &resp)
	if resp["type"] != "stream_end" {
		t.Fatalf("expected stream_end only, got: %v", resp["type"])
	}
	pl, _ := resp["payload"].(map[string]interface{})
	if pl == nil {
		t.Fatalf("expected payload")
	}
	if pl["content"] != "summary" {
		t.Fatalf("expected summary, got: %v", pl["content"])
	}

	logPath := filepath.Join(logDir, "user_a", "sess-1.log")
	raw, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log error: %v", err)
	}
	logStr := string(raw)
	if !strings.Contains(logStr, "[RAW_FINAL] raw") {
		t.Fatalf("expected raw final saved in log")
	}
	if !strings.Contains(logStr, "[SUMMARY] summary") {
		t.Fatalf("expected summary saved in log")
	}
	if !strings.Contains(logStr, "[THINKING] abc") {
		t.Fatalf("expected thinking saved in log")
	}
}
