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

	// Merchant
	http.HandleFunc("/", handlers.RenderLogin(templatesDir))
	http.HandleFunc("/api/login", handlers.HandleLogin)

	// Admin
	http.HandleFunc("/admin", handlers.RenderAdminLogin(templatesDir))
	http.HandleFunc("/api/admin/login", handlers.AdminLogin)

	http.HandleFunc("/admin/home", middlewares.Authenticate(handlers.RenderAdmin(templatesDir)))
	http.HandleFunc("/api/admin/onboard", middlewares.Authenticate(handlers.HandleOnboardUser))

	http.HandleFunc("/api/admin/users", middlewares.Authenticate(handlers.ListUsers))
	http.HandleFunc("/api/admin/edit-merchant", middlewares.Authenticate(handlers.EditMerchant))
	http.HandleFunc("/api/admin/merchants/search", middlewares.Authenticate(handlers.SearchMerchants))

	http.HandleFunc("/api/admin/transaction-logs", middlewares.Authenticate(handlers.GetTransactionLogs))
	http.HandleFunc("/api/admin/transaction-logs/search", middlewares.Authenticate(handlers.SearchTransactionLogs))

	http.HandleFunc("/api/admin/transaction-count", middlewares.Authenticate(handlers.GetSuccessTransactionCount))
	http.HandleFunc("/api/admin/total-transaction-amount", middlewares.Authenticate(handlers.GetTotalTransactionAmount))

	// Other public routes
	http.HandleFunc("/form", handlers.RenderForm(templatesDir)) // Form submission
	http.HandleFunc("/api/points", handlers.HandlePoints)       // Points management

	http.HandleFunc("/404", handlers.NotFoundHandler(templatesDir))

}
