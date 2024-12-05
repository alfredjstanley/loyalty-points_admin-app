package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

func RenderAdminLogin(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(templatesDir + "/admin-login.html")
		if err != nil {
			http.Error(w, "Unable to load admin login template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hardcoded admin credentials (you can replace this with a database query)
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if req.Username != adminUsername || req.Password != adminPassword {
		http.Error(w, `{"message":"Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	secret := os.Getenv("JWT_SECRET_ADMIN")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AdminLoginResponse{Token: tokenString})
}
