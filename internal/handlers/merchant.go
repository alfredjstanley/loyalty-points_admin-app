package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wac-offline-payment/internal/models"
	"wac-offline-payment/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type EditMerchantRequest struct {
	ID        string `json:"id"`         // Merchant ID
	StoreName string `json:"store_name"` // Updated store name
	Location  string `json:"location"`   // Updated location
	Password  string `json:"password"`   // Updated password (optional)
}

type MerchantDetailsResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func EditMerchant(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Handle View Merchant (Fetch details)
		handleViewMerchant(w, r)
	case http.MethodPut:
		// Handle Edit Merchant
		handleEditMerchant(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Function to handle the GET request for viewing merchant details
func handleViewMerchant(w http.ResponseWriter, r *http.Request) {
	// Get mobile number from query parameters
	mobileNumber := r.URL.Query().Get("mobileNumber")
	if mobileNumber == "" {
		http.Error(w, `{"success": false, "message": "Mobile number is required"}`, http.StatusBadRequest)
		return
	}

	// Step 1: Authenticate and get the token (reused from payment)
	token, err := authenticate()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Authentication failed: %s"}`, err.Error()), http.StatusUnauthorized)
		return
	}

	// Step 2: Fetch merchant details from the third-party API
	apiURL := fmt.Sprintf("https://olopo-dev.webc.in/api/merchant/details?mobileNumber=%s", mobileNumber)
	req, _ := http.NewRequest(http.MethodGet, apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Failed to fetch merchant details: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Third-party API request failed with status %d"}`, resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var result MerchantDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Error decoding response: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if !result.Success {
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "%s"}`, result.Message), http.StatusInternalServerError)
		return
	}

	// Step 3: Return the merchant details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Function to handle the PUT request for editing merchant details
func handleEditMerchant(w http.ResponseWriter, r *http.Request) {
	var req EditMerchantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate ID
	objectID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		http.Error(w, "Invalid merchant ID", http.StatusBadRequest)
		return
	}

	// Hash the new password if provided
	var hashedPassword string
	if req.Password != "" {
		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		hashedPassword = string(hashedPasswordBytes)
	}

	// Prepare update fields
	update := models.MerchantUpdate{
		StoreName: req.StoreName,
		Location:  req.Location,
		Password:  hashedPassword,
	}

	// Update merchant in the database
	if err := repository.UpdateMerchant(objectID, update); err != nil {
		http.Error(w, "Failed to update merchant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Merchant updated successfully"}`))
}
