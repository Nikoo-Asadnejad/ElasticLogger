package main

import (
	"encoding/json"
	"fmt"
	"log"
	"logger-service/logger"
	"net/http"
)

var logService *logger.Logger

func main() {

	var err error
	logService, err = logger.NewLogger("http://localhost:9200", "app-logs")
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}

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
