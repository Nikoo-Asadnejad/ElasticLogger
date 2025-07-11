package logger

import (
	"time"
)

type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Context   map[string]string `json:"context,omitempty"`
	TraceId   string            `json:"trace_id,omitempty"`
	SpanId    string            `json:"span_id,omitempty"`
	LogId     string            `json:"log_id,omitempty"`
	Type      string            `json:"type,omitempty"`
}

type ILogger interface {
	Log(entry LogEntry) error
}
