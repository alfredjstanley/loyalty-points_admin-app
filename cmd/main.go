package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"wac-offline-payment/internal/handlers"
	"wac-offline-payment/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// MongoDB URI from environment
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in environment")
	}

	// Initialize MongoDB connection
	repository.InitMongo(mongoURI)

	// Load templates directory
	templatesDir := filepath.Join("templates")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Public routes
	http.Handle("/", handlers.RenderLogin(templatesDir))
	http.HandleFunc("/api/login", handlers.HandleLogin)
	http.HandleFunc("/admin", handlers.RenderAdmin(templatesDir))
	http.HandleFunc("/api/admin/onboard", handlers.HandleOnboardUser)
	http.HandleFunc("/form", handlers.RenderForm(templatesDir))
	http.HandleFunc("/api/points", handlers.HandlePoints)
	http.HandleFunc("/api/admin/users", handlers.ListUsers)

	// Protected routes with middleware
	// protectedRoutes := http.NewServeMux()
	// // protectedRoutes.HandleFunc("/form", handlers.RenderForm(templatesDir))
	// protectedRoutes.HandleFunc("/api/points", handlers.HandlePoints)

	// http.Handle("/protected/", middleware.ValidateJWT(protectedRoutes))

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default to port 8080
	}
	log.Printf("Server started at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
