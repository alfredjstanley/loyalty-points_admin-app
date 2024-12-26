package handlers

import (
	"encoding/json"
	"net/http"

	"wac-offline-payment/internal/models"
	"wac-offline-payment/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// AddCounterHandler handles the addition of a new counter
func AddCounterHandler(w http.ResponseWriter, r *http.Request) {
	var counter models.Counter

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&counter); err != nil {
		http.Error(w, `{"success": false, "message": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if counter.MerchantPhone == "" || counter.Name == "" || counter.Location == "" || counter.Username == "" || counter.Password == "" {
		http.Error(w, `{"success": false, "message": "Missing required fields"}`, http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(counter.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	counter.Password = string(hashedPassword)

	// Add the counter
	err = repository.AddCounter(counter)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Failed to add counter"}`, http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"success": true, "message": "Counter added successfully"}`))
}

// GetCountersHandler handles fetching counters for a specific merchant
func GetCountersHandler(w http.ResponseWriter, r *http.Request) {
	merchantPhone := r.URL.Query().Get("merchant")
	if merchantPhone == "" {
		http.Error(w, `{"success": false, "message": "Merchant phone is required"}`, http.StatusBadRequest)
		return
	}

	// Fetch counters from the repository
	counters, err := repository.GetCountersByMerchant(merchantPhone)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Failed to fetch counters"}`, http.StatusInternalServerError)
		return
	}

	// Success response
	response := map[string]interface{}{
		"success":  true,
		"counters": counters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
