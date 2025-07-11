package api

import (
	"encoding/json"
	"fmt"
	"logger-service/logger"
	"net/http"
)

func HandleLog(logService logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var entry logger.LogEntry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := logService.Log(entry); err != nil {
			http.Error(w, "Failed to log: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Log entry created")
	}
}
