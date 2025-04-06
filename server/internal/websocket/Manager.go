package websocket

import (
	"context"
	"github.com/fromscript/hush/internal/websocket/models"
	"net/http"
)

type Manager interface {
	newWebSocketManager()
	UpgradeHandler(w http.ResponseWriter, r *http.Request)
	handleConnection(client *models.Client)
	readPump(ctx context.Context, c *models.Client)
	writePump(ctx context.Context, c *models.Client)
	cleanupClient(c *models.Client)
}

//type WebSocketManager struct {
//	clients          sync.Map // map[string]*Client (sessionID -> Client)
//	rooms            sync.Map // map[string]*Room (roomID -> Room)
//	msgBuffer        chan Message
//	metricsCollector metrics.Collector
//	shutdownCtx      context.Context
//	cancelFunc       context.CancelFunc
//	roomExpiry       time.Duration
//	allowedOrigins   []string
//}
