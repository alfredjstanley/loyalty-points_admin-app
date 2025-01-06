package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
	"time"

	"wac-offline-payment/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	MobileNumber string `json:"mobile_number"`
	Password     string `json:"password"`
}

type LoginResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Token       string `json:"token"`
	StoreName   string `json:"store_name"`
	Location    string `json:"location" bson:"location"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

// JWTClaims defines the structure of the JWT claims
type JWTClaims struct {
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

func RenderLogin(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(templatesDir + "/login.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming login request
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Attempt to find a user or counter
	user, counter, err := repository.FindUserOrCounter(loginReq.MobileNumber)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	var hashedPassword, phoneNumber, storeName, location string

	if user != nil {
		// User found
		hashedPassword = user.Password
		phoneNumber = user.PhoneNumber
		storeName = user.StoreName
		location = user.Location
	} else if counter != nil {
		// Counter found
		hashedPassword = counter.Password
		phoneNumber = counter.MerchantPhone
		storeName = counter.Name
		location = counter.Location
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := generateJWT(phoneNumber)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the generated token and additional information
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success:     true,
		Message:     "Login successful",
		Token:       token,
		StoreName:   storeName,
		Location:    location,
		PhoneNumber: phoneNumber,
	})
}

// generateJWT creates a JWT token for a given phone number
func generateJWT(phoneNumber string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET is not set in environment variables")
	}

	// Token expiry (default: 1 hour)
	tokenExpiry := time.Hour
	if expiry, err := time.ParseDuration(os.Getenv("TOKEN_EXPIRY") + "s"); err == nil {
		tokenExpiry = expiry
	}

	claims := JWTClaims{
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
