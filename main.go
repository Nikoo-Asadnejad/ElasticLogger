package main

import (
	"log"
	"logger-service/logger"
)

func main() {
	elasticURL := "http://localhost:9200"
	indexName := "app-logs"

	loggerService, err := logger.NewLogger(elasticURL, indexName)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	err = loggerService.Log("INFO", "Application started", map[string]string{
		"module": "main",
		"env":    "dev",
	})
	if err != nil {
		log.Println("Failed to log:", err)
	}
}
