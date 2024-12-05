package routes

import (
	"net/http"

	"wac-offline-payment/internal/handlers"
	"wac-offline-payment/internal/middlewares"
)

// RegisterRoutes registers all application routes.
func RegisterRoutes(templatesDir string) {
	// Static file server for CSS, JS, and assets
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Public routes
	http.HandleFunc("/", handlers.RenderLogin(templatesDir)) // Login page
	http.HandleFunc("/api/login", handlers.HandleLogin)      // Login API
	http.HandleFunc("/404", handlers.NotFoundHandler(templatesDir))
	// Protected routes (admin-related functionality)
	http.HandleFunc("/admin", handlers.RenderAdminLogin(templatesDir))
	http.HandleFunc("/api/admin/login", handlers.AdminLogin)

	http.HandleFunc("/admin/home", middlewares.Authenticate(handlers.RenderAdmin(templatesDir))) // Admin dashboard
	http.HandleFunc("/api/admin/onboard", middlewares.Authenticate(handlers.HandleOnboardUser))
	http.HandleFunc("/api/admin/edit-merchant", middlewares.Authenticate(handlers.EditMerchant))
	http.HandleFunc("/api/admin/merchants/search", middlewares.Authenticate(handlers.SearchMerchants))
	http.HandleFunc("/api/admin/transaction-logs", middlewares.Authenticate(handlers.GetTransactionLogs))
	http.HandleFunc("/api/admin/users", middlewares.Authenticate(handlers.ListUsers))

	// Other public routes
	http.HandleFunc("/form", handlers.RenderForm(templatesDir)) // Form submission
	http.HandleFunc("/api/points", handlers.HandlePoints)       // Points management

}
