package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"wac-offline-payment/internal/models"
	"wac-offline-payment/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type OnboardRequest struct {
	StoreName   string `json:"store_name"`
	Location    string `json:"location"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func RenderAdmin(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(templatesDir + "/admin.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func HandleOnboardUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req OnboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Save user to MongoDB
	user := models.User{
		StoreName:   req.StoreName,
		Location:    req.Location,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
	}
	if err := repository.SaveUser(user); err != nil {
		log.Printf("Error saving user: %v", err)
		http.Error(w, "Failed to onboard user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"User onboarded successfully"}`))
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default to 10 users per page
	}

	// Fetch users with pagination
	users, total, err := repository.GetUsersWithPagination(page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"users":      users,
		"total":      total, // Total merchants
		"page":       page,
		"limit":      limit,
		"totalPages": (total + limit - 1) / limit, // Calculate total pages
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type EditMerchantRequest struct {
	ID        string `json:"id"`         // Merchant ID
	StoreName string `json:"store_name"` // Updated store name
	Location  string `json:"location"`   // Updated location
	Password  string `json:"password"`   // Updated password
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

	// Update merchant details
	update := models.MerchantUpdate{
		StoreName: req.StoreName,
		Location:  req.Location,
		Password:  hashedPassword,
	}
	if err := repository.UpdateMerchant(objectID, update); err != nil {
		http.Error(w, "Failed to update merchant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Merchant updated successfully"}`))
}
