package metrics

import (
	"log/slog"
	"time"
)

type Collector interface {
	IncrementConnection()
	DecrementConnection()
	RecordMessageReceived()
	RecordMessageSent()
	RecordAuthFailure()
	RecordUpgradeFailure()
	RecordLatency(duration time.Duration)
	RecordPing(s string)
	RecordPong(s string)
}

type DefaultCollector struct{}

func (mc *DefaultCollector) IncrementConnection() {
	slog.Info("New connection established")
}

func (mc *DefaultCollector) DecrementConnection() {
	slog.Info("Connection closed")
}

func (mc *DefaultCollector) RecordMessageReceived() {
	slog.Info("New message received")
}

func (mc *DefaultCollector) RecordMessageSent() {
	slog.Info("New message sent")
}

func (mc *DefaultCollector) RecordAuthFailure() {
	slog.Info("Authentication failure")
}

func (mc *DefaultCollector) RecordUpgradeFailure() {
	slog.Info("Upgrade failure")
}
func (mc *DefaultCollector) RecordPing(s string) {
	slog.Info("New ping", s)
}

func (mc *DefaultCollector) RecordPong(s string) {
	slog.Info("New pong", s)
}

func (mc *DefaultCollector) RecordLatency(duration time.Duration) {
	slog.Info("New latency", duration)
}
