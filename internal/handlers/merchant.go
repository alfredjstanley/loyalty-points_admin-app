package handlers

import (
	"encoding/json"
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

func EditMerchant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

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
