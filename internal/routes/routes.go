package routes

import (
	"net/http"

	"wac-offline-payment/internal/handlers"
)

// RegisterRoutes registers all application routes.
func RegisterRoutes(templatesDir string) {
	// Static file server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Public routes
	http.Handle("/", handlers.RenderLogin(templatesDir))
	http.HandleFunc("/api/login", handlers.HandleLogin)
	http.HandleFunc("/admin", handlers.RenderAdmin(templatesDir))
	http.HandleFunc("/api/admin/onboard", handlers.HandleOnboardUser)
	http.HandleFunc("/form", handlers.RenderForm(templatesDir))
	http.HandleFunc("/api/points", handlers.HandlePoints)
	http.HandleFunc("/api/admin/users", handlers.ListUsers)
}
