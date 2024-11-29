package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wac-offline-payment/internal/repository"
)

func GetTransactionLogs(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Fetch logs with pagination
	logs, total, err := repository.GetLogsWithPagination(page, limit)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Failed to fetch logs"}`, http.StatusInternalServerError)
		return
	}

	// Respond with logs and metadata
	response := map[string]interface{}{
		"success": true,
		"logs":    logs,
		"total":   total,
		"limit":   limit,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
