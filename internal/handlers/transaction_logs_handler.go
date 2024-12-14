package handlers

import (
	"encoding/json"
	"fmt"
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

func GetSuccessTransactionCount(w http.ResponseWriter, r *http.Request) {
	count, err := repository.GetSuccessTransactionCount()
	if err != nil {
		http.Error(w, `{"success": false, "message": "Failed to fetch transaction count"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"count":   count,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetTotalTransactionAmount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello from 'GetTotalTransactionAmount'")
	totalAmount, err := repository.GetTotalTransactionAmount()
	if err != nil {
		http.Error(w, `{"success": false, "message": "Failed to fetch total transaction amount"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"totalAmount": totalAmount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
