package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fromscript/hush/internal/crypto"
	"github.com/fromscript/hush/internal/database"
	"github.com/fromscript/hush/internal/websocket"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Simple metrics collector implementation
type MetricsCollector struct{}

func (mc *MetricsCollector) IncrementConnection() {
	log.Println("New connection established")
}

func (mc *MetricsCollector) DecrementConnection() {
	log.Println("Connection closed")
}

func (mc *MetricsCollector) RecordMessageReceived() {}

func (mc *MetricsCollector) RecordMessageSent() {}

func (mc *MetricsCollector) RecordLatency(duration time.Duration) {}

// JWT Authenticator implementation
type JwtAuthenticator struct {
	secret []byte
}

func (a *JwtAuthenticator) ValidateToken(token string) (string, error) {
	// Implement your actual JWT validation logic here
	return "anonymous-user", nil // Demo implementation
}

func main() {
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Generate/load master encryption key
	masterKey, _ := crypto.GenerateMasterKey() // Implement secure key storage in production!

	// 3. Initialize database
	database.InitDB()
	defer database.DB.Close()

	// 4. Initialize WebSocket manager
	auth := &JwtAuthenticator{secret: []byte(os.Getenv("JWT_SECRET"))}
	metrics := &MetricsCollector{}
	wsManager := websocket.NewWebSocketManager(
		auth,
		metrics,
		masterKey,
	)

	// 5. Configure HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/ws", wsManager.UpgradeHandler)
	router.HandleFunc("/health", healthCheckHandler)

	// 6. Add middleware
	router.Use(securityHeadersMiddleware)
	router.Use(rateLimitMiddleware)

	// 7. Configure HTTP server
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 8. Graceful shutdown setup
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 9. Wait for shutdown signal
	<-done
	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wsManager.Shutdown()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement your rate limiting logic here
		next.ServeHTTP(w, r)
	})
}
