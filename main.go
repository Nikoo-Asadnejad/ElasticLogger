package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"logger-service/logger"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var logService logger.ILogger

func main() {
	loadConfig()

	var err error
	logService, err = logger.NewElasticLogger(
		viper.GetString("elasticsearch.url"),
		viper.GetString("elasticsearch.username"),
		viper.GetString("elasticsearch.password"),
		viper.GetString("elasticsearch.index"),
	)
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

	var entry logger.LogEntry
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
	conn, err := amqp.Dial(viper.GetString("rabbitmq.connection_string"))
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
		var entry logger.LogEntry
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

func loadConfig() {
	viper.SetConfigName("config") // Name of the config file (without extension)
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(".")      // Path to look for the config file

	// Read environment variables
	viper.AutomaticEnv()

	// Load the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}
