package models

import (
	"encoding/json"
)

type Message struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"timestamp,omitempty"`
	SessionID string          `json:"sessionId,omitempty"`
}
