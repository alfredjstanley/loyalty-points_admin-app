package handlers

import (
	"encoding/json"
	"net/http"

	"wac-offline-payment/internal/repository"
)

// GetMerchantReportsHandler serves the aggregated merchant reports
func GetMerchantReportsHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch reports from the repository
	reports, err := repository.GetMerchantReports()
	if err != nil {
		http.Error(w, `{"success": false, "message": "Failed to fetch reports"}`, http.StatusInternalServerError)
		return
	}

	// Build the response
	response := map[string]interface{}{
		"success": true,
		"reports": reports,
		"total":   len(reports),
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
