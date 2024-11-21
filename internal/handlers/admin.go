package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"wac-offline-payment/internal/models"
	"wac-offline-payment/internal/repository"

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

	users, err := repository.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
