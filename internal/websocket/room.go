package websocket

import (
	"sync"
	"time"
)

// Room represents a chat room with members and activity tracking
type Room struct {
	ID           string
	members      sync.Map // map[string]struct{} (session IDs of members)
	createdAt    time.Time
	lastActivity time.Time
	mu           sync.RWMutex // Protects activity timestamps
}

// NewRoom creates a new Room instance with proper initialization
func NewRoom(id string) *Room {
	now := time.Now()
	return &Room{
		ID:           id,
		createdAt:    now,
		lastActivity: now,
	}
}

// AddMember adds a client to the room
func (r *Room) AddMember(sessionID string) {
	r.members.Store(sessionID, struct{}{})
	r.updateActivity()
}

// RemoveMember removes a client from the room
func (r *Room) RemoveMember(sessionID string) {
	r.members.Delete(sessionID)
	r.updateActivity()
}

// GetMemberCount returns the number of active members
func (r *Room) GetMemberCount() int {
	count := 0
	r.members.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// GetActiveDuration returns time since last activity
func (r *Room) GetActiveDuration() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return time.Since(r.lastActivity)
}

// updateActivity marks the room as active
func (r *Room) updateActivity() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastActivity = time.Now()
}

// Expired checks if the room has been inactive beyond threshold
func (r *Room) Expired(expiryThreshold time.Duration) bool {
	return r.GetActiveDuration() > expiryThreshold
}

// Cleanup performs graceful shutdown of the room
func (r *Room) Cleanup(wm *WebSocketManager) {
	r.members.Range(func(sessionID, _ interface{}) bool {
		if client, ok := wm.clients.Load(sessionID); ok {
			client.(*Client).roomID = ""
		}
		return true
	})
}
