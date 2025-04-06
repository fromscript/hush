package models

import (
	"github.com/coder/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	SessionID string
	Send      chan Message
	RoomID    string
}
