package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretkey")

// Context key for storing user email
type contextKey string

const UserEmailKey contextKey = "user_email"

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Expect header like: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["email"] == nil {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add email to context
		ctx := context.WithValue(r.Context(), UserEmailKey, claims["email"].(string))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}



/*
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        email := "bhavya@example.com" // pretend extracted from token
        ctx := context.WithValue(r.Context(), "user_email", email)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
    email := r.Context().Value("user_email").(string)
    fmt.Fprintf(w, "Hello, %s", email)
}
*/