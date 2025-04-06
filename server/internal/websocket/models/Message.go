package models

import "encoding/json"

type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	SessionID string          `json:"sessionId,omitempty"`
}
