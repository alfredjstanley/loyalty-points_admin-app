package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"
	"wac-offline-payment/internal/handlers"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserContextKey ContextKey = "user"

// ValidateJWT middleware validates the JWT token
func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			http.Error(w, "Server error: JWT_SECRET not configured", http.StatusInternalServerError)
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &handlers.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*handlers.JWTClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user phone number to request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims.PhoneNumber)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
