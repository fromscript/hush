package model

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	writeTimeout   = 10 * time.Second
	pingInterval   = 30 * time.Second
	maxMessageSize = 1024 * 1024 // 1MB
)

type WebSocketManager struct {
	clients   sync.Map // map[string]*Client
	rooms     sync.Map // map[string]*Room
	msgBuffer chan Message
}

type Client struct {
	conn      *websocket.Conn
	sessionID string
	userID    string
	roomID    string
	send      chan Message
}

type Message struct {
	Type    string          `json:"type"` // message, control, auth
	Payload json.RawMessage `json:"payload"`
	Session string          `json:"session,omitempty"`
}
