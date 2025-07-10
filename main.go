package main

import (
	"encoding/json"
	"fmt"
	"log"
	"logger-service/elasticlogger"
	"net/http"

	"github.com/streadway/amqp"
)

var logService elasticlogger.ILogger

func main() {

	var err error
	logService, err = elasticlogger.NewElasticLogger("http://localhost:9200", "elastic", "1234", "app-logs")
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}

	go startRabbitMQConsumer()

	http.HandleFunc("/log", handleLog)

	fmt.Println("Listening on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleLog(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry elasticlogger.LogEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = logService.Log(entry)
	if err != nil {
		http.Error(w, "Failed to log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Log entry created")
}

func startRabbitMQConsumer() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"log_entries", // queue name
		true,          // durable
		false,         // auto-delete
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer tag
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	fmt.Println("RabbitMQ consumer started. Waiting for messages...")
	for msg := range msgs {
		var entry elasticlogger.LogEntry
		err := json.Unmarshal(msg.Body, &entry)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		err = logService.Log(entry)
		if err != nil {
			log.Printf("Failed to log entry: %v", err)
		} else {
			fmt.Println("Log entry successfully logged:", entry)
		}
	}
}
