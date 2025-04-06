package handlers

import (
	"github.com/fromscript/hush/internal/crypto"
	"github.com/fromscript/hush/internal/database"
	_ "log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID string
	RoomID string
	Send   chan []byte
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan Message)
)

func HandleWebSocket(conn *websocket.Conn, roomKey []byte) {
	client := &Client{Conn: conn, Send: make(chan []byte)}
	clients[client] = true

	defer func() {
		delete(clients, client)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Decrypt message with room key
		decrypted, _ := crypto.Decrypt(msg, roomKey)

		// Process message (store in DB, broadcast)
		database.SaveMessage(string(decrypted))
		broadcast <- Message{Content: msg, RoomID: client.RoomID}
	}
}

func WebSocketHandler(key []byte) func(http.ResponseWriter, *http.Request) {

}
