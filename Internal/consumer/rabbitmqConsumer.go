package rabbitMqConsumer

import (
	"encoding/json"
	"fmt"
	"log"
	"logger-service/internal/logger"

	"github.com/streadway/amqp"
)

func StartConsumer(logService logger.ILogger, connectionString string) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("log_entries", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	fmt.Println("RabbitMQ consumer started. Waiting for messages...")
	for msg := range msgs {
		var entry logger.LogEntry
		if err := json.Unmarshal(msg.Body, &entry); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		if err := logService.Log(entry); err != nil {
			log.Printf("Failed to log entry: %v", err)
		} else {
			fmt.Println("Log entry successfully logged:", entry)
		}
	}
}
