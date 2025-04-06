package websocket

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/fromscript/hush/internal/websocket/models"
)

const (
	writeTimeout   = 10 * time.Second
	pingInterval   = 30 * time.Second
	maxMessageSize = 1024 * 1024 // 1MB
	NormalClosure  = websocket.StatusNormalClosure
	InternalError  = websocket.StatusInternalError
)

type DefaultManager struct {
	clients   sync.Map // map[string]*models.Client
	rooms     sync.Map // map[string]*models.Room
	authToken string
}

func NewDefaultManager(authToken string) *DefaultManager {
	return &DefaultManager{
		authToken: authToken,
	}
}

func (dm *DefaultManager) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("token") != dm.authToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}

	sessionID, _ := generateSessionID()
	client := &models.Client{
		Conn:      conn,
		SessionID: sessionID,
		Send:      make(chan models.Message, 256),
	}

	dm.clients.Store(sessionID, client)
	go dm.handleConnection(client)
}

func (dm *DefaultManager) handleConnection(client *models.Client) {
	defer dm.cleanupClient(client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go dm.readPump(ctx, client)
	dm.writePump(ctx, client)
}

func (dm *DefaultManager) readPump(ctx context.Context, client *models.Client) {
	client.Conn.SetReadLimit(maxMessageSize)

	for {
		_, data, err := client.Conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				slog.Info("Client disconnected", "session", client.SessionID)
			} else {
				slog.Warn("Read error", "session", client.SessionID, "error", err)
			}
			return
		}

		var msg models.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			slog.Warn("Invalid message format", "session", client.SessionID, "error", err)
			continue
		}

		dm.processMessage(client, msg)
	}
}

func (dm *DefaultManager) processMessage(client *models.Client, msg models.Message) {
	switch msg.Type {
	case "join":
		var joinMsg models.JoinMessage
		if err := json.Unmarshal(msg.Payload, &joinMsg); err == nil {
			dm.joinRoom(client, joinMsg.RoomID)
			dm.sendSystemMessage(client, "joined", joinMsg.RoomID)
		}
	case "message":
		if client.RoomID != "" {
			dm.broadcastToRoom(client.RoomID, msg)
		}
	default:
		slog.Warn("Unknown message type", "type", msg.Type)
	}
}

func (dm *DefaultManager) writePump(ctx context.Context, client *models.Client) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case msg := <-client.Send:
			err := wsjson.Write(ctx, client.Conn, msg)
			if err != nil {
				slog.Warn("Write error", "session", client.SessionID, "error", err)
				return
			}

		case <-ticker.C:
			if err := client.Conn.Ping(ctx); err != nil {
				slog.Warn("Ping failed", "session", client.SessionID, "error", err)
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (dm *DefaultManager) getOrCreateRoom(roomID string) *models.Room {
	actual, _ := dm.rooms.LoadOrStore(roomID, &models.Room{
		ID:      roomID,
		Members: sync.Map{},
	})
	return actual.(*models.Room)
}

func (dm *DefaultManager) joinRoom(client *models.Client, roomID string) {
	// Leave previous room
	if client.RoomID != "" {
		if room, ok := dm.rooms.Load(client.RoomID); ok {
			room.(*models.Room).Members.Delete(client.SessionID)
		}
	}

	room := dm.getOrCreateRoom(roomID)
	room.Members.Store(client.SessionID, client)
	client.RoomID = roomID
	slog.Info("Client joined room", "session", client.SessionID, "room", roomID)
}

func (dm *DefaultManager) broadcastToRoom(roomID string, msg models.Message) {
	if room, ok := dm.rooms.Load(roomID); ok {
		room.(*models.Room).Members.Range(func(_, value interface{}) bool {
			client := value.(*models.Client)
			if client.SessionID != msg.SessionID {
				select {
				case client.Send <- msg:
				default:
					slog.Warn("Client buffer full", "session", client.SessionID)
				}
			}
			return true
		})
	}
}

func (dm *DefaultManager) sendSystemMessage(client *models.Client, msgType string, data interface{}) {
	payload, _ := json.Marshal(data)
	client.Send <- models.Message{
		Type:    "system",
		Payload: payload,
	}
}

func (dm *DefaultManager) cleanupClient(client *models.Client) {
	dm.clients.Delete(client.SessionID)
	client.Conn.Close(NormalClosure, "Connection closed")
	close(client.Send)

	if client.RoomID != "" {
		if room, ok := dm.rooms.Load(client.RoomID); ok {
			room.(*models.Room).Members.Delete(client.SessionID)
		}
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
