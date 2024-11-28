package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"wac-offline-payment/internal/repository"
	"wac-offline-payment/internal/routes"

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

	// Register routes
	routes.RegisterRoutes(templatesDir)

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
