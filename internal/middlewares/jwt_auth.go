package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Authenticate checks if the user is authenticated by verifying the JWT token.
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is correct
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET_ADMIN")), nil
		})

		fmt.Print(token)

		if err != nil || !token.Valid {
			fmt.Printf(err.Error())
			http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
			return
		}

		next(w, r)
	}
}
