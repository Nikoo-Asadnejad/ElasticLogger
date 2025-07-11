package main

import (
	"fmt"
	"log"
	"logger-service/Internal/config"
	rabbitMqConsumer "logger-service/internal/consumer"
	"net/http"

	"logger-service/internal/api"
	"logger-service/internal/logger"

	"github.com/spf13/viper"
)

var logService logger.ILogger

func main() {

	config.Load()

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

	go rabbitMqConsumer.StartConsumer(logService, viper.GetString("rabbitmq.connection_string"))

	http.HandleFunc("/log", api.HandleLog(logService))
	fmt.Println("Listening on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
