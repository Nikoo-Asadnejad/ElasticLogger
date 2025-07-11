package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticLogger struct {
	elasticClient *elasticsearch.Client
	index         string
}

func NewElasticLogger(elasticUrl, username, password, index string) (*ElasticLogger, error) {

	elasticConfig := elasticsearch.Config{
		Addresses: []string{elasticUrl},
		Password:  password,
		Username:  username,
	}

	elasticClient, err := elasticsearch.NewClient(elasticConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	return &ElasticLogger{elasticClient: elasticClient, index: index}, nil
}

func (logger *ElasticLogger) Log(entry LogEntry) error {

	entry.Timestamp = time.Now()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	res, err := logger.elasticClient.Index(
		logger.index,
		bytes.NewReader(data),
		logger.elasticClient.Index.WithContext(context.Background()),
		logger.elasticClient.Index.WithRefresh("true"),
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
