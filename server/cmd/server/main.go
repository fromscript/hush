package main

import (
	"encoding/json"
	"github.com/fromscript/hush/internal/websocket"
	"log"
	"net/http"
	"time"
)

func main() {
	manager := websocket.NewDefaultManager("development-token")

	http.HandleFunc("/ws", manager.UpgradeHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"version": "1.0.0",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
