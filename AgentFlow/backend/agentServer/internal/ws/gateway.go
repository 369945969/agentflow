package ws

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"agentServer/internal/agentcore"
	"agentServer/internal/database"
	"agentServer/internal/protocol"
)

type MessageHandler func(ctx context.Context, connID string, msg protocol.Message)

type Gateway struct {
	upgrader websocket.Upgrader
	handler  MessageHandler
	core     *agentcore.AgentCore // Reference to AgentCore for services

	mu          sync.RWMutex
	conns       map[string]*Conn
	sessions    map[string]map[string]struct{}
	connSession map[string]string

	pingInterval  time.Duration
	pongWait      time.Duration
	sendQueueSize int

	onceStart sync.Once
}

type Option func(*Gateway)

func WithHeartbeat(pingInterval time.Duration, pongWait time.Duration) Option {
	return func(g *Gateway) {
		if pingInterval > 0 {
			g.pingInterval = pingInterval
		}
		if pongWait > 0 {
			g.pongWait = pongWait
		}
	}
}

func WithSendQueueSize(size int) Option {
	return func(g *Gateway) {
		if size > 0 {
			g.sendQueueSize = size
		}
	}
}

func NewGateway(upgrader websocket.Upgrader, core *agentcore.AgentCore, handler MessageHandler, opts ...Option) *Gateway {
	g := &Gateway{
		upgrader:      upgrader,
		core:          core,
		handler:       handler,
		conns:         make(map[string]*Conn),
		sessions:      make(map[string]map[string]struct{}),
		connSession:   make(map[string]string),
		pingInterval:  30 * time.Second,
		pongWait:      60 * time.Second,
		sendQueueSize: 100,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.onceStart.Do(func() {
		go g.heartbeatLoop()
	})

	wsConn, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade failed remote=%s err=%v", r.RemoteAddr, err)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	connID := generateID()
	log.Printf("[ws] connection established conn_id=%s remote=%s path=%s", connID, r.RemoteAddr, r.URL.Path)
	conn := newConn(
		connID,
		wsConn,
		g,
		g.sendQueueSize,
		ctx.Done(),
		cancel,
		func(id string) {
			log.Printf("[ws] connection closed conn_id=%s", id)
			g.mu.Lock()
			delete(g.conns, id)
			if sid, ok := g.connSession[id]; ok {
				delete(g.connSession, id)
				if m, ok := g.sessions[sid]; ok {
					delete(m, id)
					if len(m) == 0 {
						delete(g.sessions, sid)
					}
				}
			}
			g.mu.Unlock()
		},
		func(id string, sessionID string) {
			log.Printf("[ws] session bound conn_id=%s session_id=%s", id, sessionID)
			g.mu.Lock()
			defer g.mu.Unlock()
			prev, ok := g.connSession[id]
			if ok && prev == sessionID {
				return
			}
			if ok && prev != "" && prev != sessionID {
				if m, ok := g.sessions[prev]; ok {
					delete(m, id)
					if len(m) == 0 {
						delete(g.sessions, prev)
					}
				}
			}
			g.connSession[id] = sessionID
			m, ok := g.sessions[sessionID]
			if !ok {
				m = make(map[string]struct{})
				g.sessions[sessionID] = m
			}
			m[id] = struct{}{}
		},
	)

	conn.wsConn.SetReadDeadline(time.Now().Add(g.pongWait))
	conn.wsConn.SetPingHandler(func(appData string) error {
		atomic.StoreInt64(&conn.lastPingUnix, time.Now().Unix())
		_ = conn.wsConn.SetReadDeadline(time.Now().Add(g.pongWait))
		_ = conn.wsConn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(5*time.Second))
		return nil
	})
	conn.wsConn.SetPongHandler(func(string) error {
		atomic.StoreInt64(&conn.lastPingUnix, time.Now().Unix())
		_ = conn.wsConn.SetReadDeadline(time.Now().Add(g.pongWait))
		return nil
	})

	g.mu.Lock()
	g.conns[connID] = conn
	g.mu.Unlock()

	go conn.writeLoop()
	conn.readLoop(ctx, g.handler)
}

type Conn struct {
	ID        string
	UserID    string
	gateway   *Gateway
	wsConn    *websocket.Conn
	sendCh    chan protocol.ServerMessage
	pingCh    chan struct{}
	done      <-chan struct{}
	onSession func(connID string, sessionID string)

	closeOnce sync.Once
	cancel    context.CancelFunc
	onClose   func(connID string)

	lastPingUnix int64
}

func newConn(
	id string,
	wsConn *websocket.Conn,
	gateway *Gateway,
	sendQueueSize int,
	done <-chan struct{},
	cancel context.CancelFunc,
	onClose func(string),
	onSession func(string, string),
) *Conn {
	return &Conn{
		ID:           id,
		wsConn:       wsConn,
		gateway:      gateway,
		sendCh:       make(chan protocol.ServerMessage, sendQueueSize),
		pingCh:       make(chan struct{}, 1),
		done:         done,
		cancel:       cancel,
		onClose:      onClose,
		onSession:    onSession,
		lastPingUnix: time.Now().Unix(),
	}
}

func (c *Conn) Send(msg protocol.ServerMessage) {
	select {
	case c.sendCh <- msg:
	default:
		c.Close()
	}
}

func (c *Conn) EnqueuePing() {
	select {
	case c.pingCh <- struct{}{}:
	default:
	}
}

func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		c.cancel()
		close(c.sendCh)
		close(c.pingCh)
		_ = c.wsConn.Close()
		if c.onClose != nil {
			c.onClose(c.ID)
		}
	})
}

func (c *Conn) readLoop(ctx context.Context, handler MessageHandler) {
	defer c.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, data, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Printf("[ws] read failed conn_id=%s err=%v", c.ID, err)
			return
		}
		atomic.StoreInt64(&c.lastPingUnix, time.Now().Unix())
		log.Printf("[ws] received raw message conn_id=%s payload=%s", c.ID, string(data))

		var msg protocol.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("[ws] invalid json conn_id=%s err=%v", c.ID, err)
			c.Send(protocol.ServerMessage{
				Type:      "error",
				Timestamp: time.Now().UnixMilli(),
				Payload: protocol.ErrorPayload{
					Code:      "invalid_json",
					Message:   "Invalid JSON format",
					Retryable: false,
				},
			})
			continue
		}

		if msg.MessageID == "" || msg.Content.Text == "" || (msg.GroupID == "" && msg.UserID == "") {
			log.Printf("[ws] validation failed conn_id=%s message_id=%s user_id=%s group_id=%s", c.ID, msg.MessageID, msg.UserID, msg.GroupID)
			c.Send(protocol.ServerMessage{
				Type:      "error",
				SessionID: msg.SessionID,
				MessageID: msg.MessageID,
				Timestamp: time.Now().UnixMilli(),
				Payload: protocol.ErrorPayload{
					Code:      "validation_error",
					Message:   "message_id, content.text and (group_id or user_id) are required",
					Retryable: false,
				},
			})
			continue
		}

		if msg.SessionID == "" {
			if msg.GroupID != "" {
				msg.SessionID = msg.GroupID
			} else {
				msg.SessionID = generateID()
			}
		}

		if c.onSession != nil {
			c.onSession(c.ID, msg.SessionID)
		}

		if msg.UserID != "" {
			if c.UserID == "" {
				c.UserID = msg.UserID
				log.Printf("[ws] user bound conn_id=%s user_id=%s", c.ID, c.UserID)
			} else if c.UserID != msg.UserID {
				log.Printf("[ws] user mismatch conn_id=%s current_user=%s incoming_user=%s", c.ID, c.UserID, msg.UserID)
				c.Send(protocol.ServerMessage{
					Type:      "error",
					SessionID: msg.SessionID,
					MessageID: msg.MessageID,
					Timestamp: time.Now().UnixMilli(),
					Payload: protocol.ErrorPayload{
						Code:      "user_mismatch",
						Message:   "user_id mismatch",
						Retryable: false,
					},
				})
				continue
			}
		}

		log.Printf("[ws] processing message conn_id=%s session_id=%s message_id=%s type=%s", c.ID, msg.SessionID, msg.MessageID, msg.Type)

		// Only process messages with user_id (single chat)
		if msg.UserID != "" {
			go c.processSingleChatMessage(ctx, msg)
		} else {
			// For group messages (if needed), but focus on single chat first
			go handler(ctx, c.ID, msg)
		}
	}
}

// processSingleChatMessage handles single chat messages by loading user settings
// and forwarding to opencode via SSE
func (c *Conn) processSingleChatMessage(ctx context.Context, msg protocol.Message) {
	// Get the Gateway reference
	g := c.gateway
	if g == nil || g.core == nil {
		log.Printf("[ws] gateway or core not initialized for conn_id=%s", c.ID)
		c.Send(protocol.ServerMessage{
			Type:      "error",
			MessageID: msg.MessageID,
			SessionID: msg.SessionID,
			Timestamp: time.Now().UnixMilli(),
			Payload: protocol.ErrorPayload{
				Code:      "system_error",
				Message:   "System not initialized",
				Retryable: false,
			},
		})
		return
	}

	// Load user settings from database
	settings, err := database.GetUserSettings(ctx, msg.UserID)
	if err != nil {
		log.Printf("[ws] failed to load user settings conn_id=%s user_id=%s err=%v", c.ID, msg.UserID, err)
		c.Send(protocol.ServerMessage{
			Type:      "error",
			MessageID: msg.MessageID,
			SessionID: msg.SessionID,
			Timestamp: time.Now().UnixMilli(),
			Payload: protocol.ErrorPayload{
				Code:      "settings_error",
				Message:   "Failed to load user settings",
				Retryable: true,
			},
		})
		return
	}

	// Add user settings to message
	msg.ThinkingEnabled = settings.ThinkingEnabled
	msg.SimplifiedOutput = settings.SimplifiedOutput
	msg.SystemPrompt = settings.SystemPrompt
	msg.BoundModel = settings.BoundModel

	log.Printf("[ws] loaded user settings conn_id=%s user_id=%s thinking=%t simplified=%t model=%s",
		c.ID, msg.UserID, msg.ThinkingEnabled, msg.SimplifiedOutput, msg.BoundModel)

	// Use HandleMessage for orchestration, it handles streaming via emitter
	g.core.HandleMessage(ctx, c.ID, msg)
}

// streamResponses streams SSE responses back to the WebSocket client
func (c *Conn) streamResponses(ctx context.Context, msg protocol.Message, out <-chan protocol.ServerMessage, errCh <-chan error) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("[ws] stream cancelled conn_id=%s session_id=%s", c.ID, msg.SessionID)
			return

		case err, ok := <-errCh:
			if ok {
				log.Printf("[ws] stream error conn_id=%s err=%v", c.ID, err)
				c.Send(protocol.ServerMessage{
					Type:      "error",
					MessageID: msg.MessageID,
					SessionID: msg.SessionID,
					Timestamp: time.Now().UnixMilli(),
					Payload: protocol.ErrorPayload{
						Code:      "stream_error",
						Message:   "Stream error: " + err.Error(),
						Retryable: false,
					},
				})
			}
			return

		case resp, ok := <-out:
			if !ok {
				return // Stream ended
			}

			// Send the response to the client
			c.Send(resp)
		}
	}
}

func (c *Conn) writeLoop() {
	defer c.Close()

	for {
		select {
		case <-c.done:
			return
		case _, ok := <-c.pingCh:
			if !ok {
				return
			}
			log.Printf("[ws] sending ping conn_id=%s", c.ID)
			_ = c.wsConn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second))
		case msg, ok := <-c.sendCh:
			if !ok {
				return
			}
			log.Printf("[ws] sending message conn_id=%s type=%s session_id=%s message_id=%s", c.ID, msg.Type, msg.SessionID, msg.MessageID)
			if err := c.wsConn.WriteJSON(msg); err != nil {
				log.Printf("[ws] send failed conn_id=%s err=%v", c.ID, err)
				return
			}
		}
	}
}

func (g *Gateway) heartbeatLoop() {
	ticker := time.NewTicker(g.pingInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Unix()

		g.mu.RLock()
		conns := make([]*Conn, 0, len(g.conns))
		for _, c := range g.conns {
			conns = append(conns, c)
		}
		g.mu.RUnlock()

		for _, c := range conns {
			last := atomic.LoadInt64(&c.lastPingUnix)
			if last > 0 && now-last > int64(g.pongWait.Seconds()) {
				c.Close()
				continue
			}
			c.EnqueuePing()
		}
	}
}

func generateID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func (g *Gateway) SendToConnection(connID string, msg protocol.ServerMessage) {
	g.mu.RLock()
	c := g.conns[connID]
	g.mu.RUnlock()
	if c == nil {
		return
	}
	c.Send(msg)
}

func (g *Gateway) BroadcastToSession(sessionID string, msg protocol.ServerMessage) {
	g.mu.RLock()
	m := g.sessions[sessionID]
	if len(m) == 0 {
		g.mu.RUnlock()
		return
	}
	conns := make([]*Conn, 0, len(m))
	for id := range m {
		if c, ok := g.conns[id]; ok {
			conns = append(conns, c)
		}
	}
	g.mu.RUnlock()

	for _, c := range conns {
		c.Send(msg)
	}
}
