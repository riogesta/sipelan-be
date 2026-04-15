package common

import (
	"context"
	"encoding/json"
	"net/http"
	"sipelan/database"
	"sipelan/models"
	"strings"
	"time"
)

// CORSMiddleware handles Cross-Origin Resource Sharing headers
// to allow the React frontend to communicate with the backend.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JSONMiddleware sets the Content-Type header to application/json.
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Error(http.StatusUnauthorized, "Authorization token not provided"))
			return
		}

		// Strip "Bearer " prefix if present (case-insensitive)
		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			tokenString = strings.TrimSpace(tokenString[7:])
		} else {
			tokenString = strings.TrimSpace(tokenString)
		}

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Error(http.StatusUnauthorized, "Invalid authorization token format"))
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Error(http.StatusUnauthorized, "Invalid authorization token: "+err.Error()))
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Error(http.StatusUnauthorized, "Token has expired"))
			return
		}

		if claims.ID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Error(http.StatusUnauthorized, "Invalid authorization token: person ID is 0"))
			return
		}

		var person models.Person
		result := database.DB.Where("id = ?", claims.ID).First(&person)
		if result.Error != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Error(http.StatusUnauthorized, "Invalid authorization token: person not found"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "person", person))
		next.ServeHTTP(w, r)
	})
}
