package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"wsServer/internal/agentcore"
	"wsServer/internal/protocol"
	wsInternal "wsServer/internal/ws"
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

	wsServer := httptest.NewServer(gateway)
	defer wsServer.Close()

	u, err := url.Parse(wsServer.URL)
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
