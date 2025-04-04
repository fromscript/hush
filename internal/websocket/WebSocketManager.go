package websocket

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/fromscript/hush/internal/crypto"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	writeTimeout        = 10 * time.Second
	pingInterval        = 30 * time.Second
	maxInactiveDuration = 4 * time.Minute
	maxMessageSize      = 1024 * 1024 // 1MB
	NormalClosure       = websocket.StatusNormalClosure
	GoingAway           = websocket.StatusGoingAway
	ProtocolError       = websocket.StatusProtocolError
	UnsupportedData     = websocket.StatusUnsupportedData
	InternalError       = websocket.StatusInternalError
)

type WebSocketManager struct {
	clients     sync.Map // map[string]*Client (sessionID -> Client)
	rooms       sync.Map // map[string]*Room (roomID -> Room)
	msgBuffer   chan Message
	masterKey   []byte
	auth        Authenticator
	metrics     MetricsCollector
	shutdownCtx context.Context
	cancelFunc  context.CancelFunc
}

type Authenticator interface {
	ValidateToken(token string) (userID string, err error)
}

type MetricsCollector interface {
	IncrementConnection()
	DecrementConnection()
	RecordMessageReceived()
	RecordMessageSent()
	RecordLatency(duration time.Duration)
}

type Client struct {
	conn         *websocket.Conn
	sessionID    string
	userID       string
	roomID       string
	send         chan Message
	lastActivity time.Time
	mu           sync.Mutex
}

type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	SessionID string          `json:"session_id,omitempty"`
	Timestamp int64           `json:"timestamp"`
}

func NewWebSocketManager(auth Authenticator, metrics MetricsCollector, masterKey []byte) *WebSocketManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketManager{
		msgBuffer:   make(chan Message, 10000),
		masterKey:   masterKey,
		auth:        auth,
		metrics:     metrics,
		shutdownCtx: ctx,
		cancelFunc:  cancel,
	}
}

func (wm *WebSocketManager) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// 1. Authentication
	token := r.URL.Query().Get("token")
	userID, err := wm.auth.ValidateToken(token)
	if err != nil {
		slog.Error("Authentication failed", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. WebSocket upgrade with security headers
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("X-Frame-Options", "DENY")

	opts := &websocket.AcceptOptions{
		CompressionMode:    websocket.CompressionContextTakeover,
		InsecureSkipVerify: false,
		OriginPatterns:     []string{"*"},
	}

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}

	// 3. Configure connection settings
	conn.SetReadLimit(maxMessageSize)
	conn.CloseRead(r.Context()) // Handle reads in separate goroutine

	// 4. Generate secure session ID
	sessionID, err := generateSessionID()
	if err != nil {
		slog.Error("Session ID generation failed", "error", err)

		conn.Close(InternalError, "Internal Server Error")
		return
	}

	client := &Client{
		conn:         conn,
		sessionID:    sessionID,
		userID:       userID,
		send:         make(chan Message, 256),
		lastActivity: time.Now(),
	}

	wm.clients.Store(sessionID, client)
	wm.metrics.IncrementConnection()

	// 5. Start processing goroutines
	go wm.readPump(r.Context(), client)
	go wm.writePump(r.Context(), client)

	wm.metrics.RecordLatency(time.Since(start))
}

func (wm *WebSocketManager) readPump(ctx context.Context, c *Client) {
	defer func() {
		wm.cleanupClient(c)
		wm.metrics.DecrementConnection()
	}()

	for {
		select {
		case <-wm.shutdownCtx.Done():
			return
		default:
			msgType, data, err := c.conn.Read(ctx)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					slog.Warn("Read error", "session", c.sessionID, "error", err)
				}
				return
			}

			if msgType != websocket.MessageText {
				slog.Warn("Invalid message type", "session", c.sessionID, "type", msgType)
				continue
			}

			// Process message in separate goroutine to prevent blocking
			go wm.processMessage(c, data)
		}
	}
}

func (wm *WebSocketManager) processMessage(c *Client, data []byte) {
	start := time.Now()
	defer wm.metrics.RecordLatency(time.Since(start))

	// 1. Decrypt message
	decrypted, err := crypto.Decrypt(data, wm.masterKey)
	if err != nil {
		slog.Error("Decryption failed", "session", c.sessionID, "error", err)
		return
	}

	// 2. Validate message structure
	var msg Message
	if err := json.Unmarshal(decrypted, &msg); err != nil {
		slog.Warn("Invalid message format", "session", c.sessionID, "error", err)
		return
	}

	// 3. Add to processing pipeline
	msg.SessionID = c.sessionID
	msg.Timestamp = time.Now().UnixNano()

	select {
	case wm.msgBuffer <- msg:
		wm.metrics.RecordMessageReceived()
	default:
		slog.Warn("Message buffer full, dropping message", "session", c.sessionID)
	}
}

func (wm *WebSocketManager) writePump(ctx context.Context, c *Client) {
	defer func() {
		wm.CloseConnection(c, NormalClosure, "connection closed")
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				wm.CloseConnection(c, NormalClosure, "server shutdown")
				return
			}

			if err := wm.sendMessage(ctx, c, msg); err != nil {
				slog.Warn("Write error", "session", c.sessionID, "error", err)
				return
			}

		case <-ticker.C:
			if err := c.conn.Ping(ctx); err != nil {
				slog.Warn("Ping failed", "session", c.sessionID, "error", err)
				return
			}
			c.updateActivity()
		}
	}
}

func (wm *WebSocketManager) sendMessage(ctx context.Context, c *Client, msg Message) error {
	start := time.Now()
	defer wm.metrics.RecordLatency(time.Since(start))

	// 1. Encrypt message
	encrypted, err := crypto.Encrypt(msg.Payload, wm.masterKey)
	if err != nil {
		return err
	}

	// 2. Prepare final message
	finalMsg := Message{
		Type:      msg.Type,
		Payload:   encrypted,
		Timestamp: time.Now().UnixNano(),
	}

	// 3. Context-aware write with timeout
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	if err := wsjson.Write(ctx, c.conn, finalMsg); err != nil {
		return err
	}

	wm.metrics.RecordMessageSent()
	return nil
}

func (wm *WebSocketManager) Broadcast(roomID string, msg Message) error {
	room, ok := wm.rooms.Load(roomID)
	if !ok {
		return errors.New("room not found")
	}

	room.(*Room).members.Range(func(key, value interface{}) bool {
		client, ok := wm.clients.Load(key.(string))
		if ok {
			select {
			case client.(*Client).send <- msg:
			default:
				slog.Warn("Client send buffer full", "session", key)
			}
		}
		return true
	})

	return nil
}

func (wm *WebSocketManager) Shutdown() {
	wm.clients.Range(func(_, value interface{}) bool {
		client := value.(*Client)
		wm.CloseConnection(client, GoingAway, "server maintenance")
		return true
	})
}

func (wm *WebSocketManager) CloseConnection(client *Client, statusCode websocket.StatusCode, reason string) {
	err := client.conn.Close(statusCode, reason)
	if err != nil {
		return
	}
}

func (c *Client) updateActivity() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastActivity = time.Now()
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (wm *WebSocketManager) cleanupClient(c *Client) {
	wm.clients.Delete(c.sessionID)
	close(c.send)
	wm.CloseConnection(c, InternalError, "connection closed")

	if c.roomID != "" {
		if room, ok := wm.rooms.Load(c.roomID); ok {
			room.(*Room).members.Delete(c.sessionID)
		}
	}
}
