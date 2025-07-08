package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Logger struct {
	es    *elasticsearch.Client
	index string
}

type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Context   map[string]string `json:"context,omitempty"`
}

// Constructor
func NewLogger(esURL, index string) (*Logger, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	return &Logger{es: es, index: index}, nil
}

func (l *Logger) Log(entry LogEntry) error {
	entry.Timestamp = time.Now()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	res, err := l.es.Index(
		l.index,
		bytes.NewReader(data),
		l.es.Index.WithContext(context.Background()),
		l.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("failed to index log: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	return nil
}
