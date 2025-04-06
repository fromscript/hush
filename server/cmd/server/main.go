package main

import (
	"github.com/fromscript/hush/internal/websocket"
	"log"
	"net/http"
)

func main() {
	manager := websocket.NewDefaultManager("development-token")

	http.HandleFunc("/ws", manager.UpgradeHandler)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
