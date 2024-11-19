package main

import (
	"log"
	"net/http"
	"path/filepath"

	"wac-offline-payment/internal/handlers"
)

func main() {
	// Load templates
	templatesDir := filepath.Join("templates")
	http.Handle("/", handlers.RenderForm(templatesDir))
	http.HandleFunc("/submit", handlers.HandlePaymentSubmission)

	// API route for points
	http.HandleFunc("/api/points", handlers.HandlePoints)

	port := "8080" // Default port for local development
	log.Printf("Server started at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
