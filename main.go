package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sipelan/common"
	"sipelan/database"
	"sipelan/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	database.Connect()

	r := chi.NewRouter()

	// Global middleware
	r.Use(common.CORSMiddleware)
	r.Use(common.JSONMiddleware)

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handlers.Auth)
	})

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(common.AuthMiddleware)

		// Categories
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", handlers.GetAllCategories)
			r.Post("/", handlers.CreateCategory)
			r.Get("/{id}", handlers.GetCategoryByID)
			r.Put("/{id}", handlers.UpdateCategory)
			r.Delete("/{id}", handlers.DeleteCategory)
		})

		// Transactions
		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", handlers.GetAllTransactions)
			r.Post("/", handlers.CreateTransaction)
			r.Get("/{id}", handlers.GetTransactionByID)
			r.Put("/{id}", handlers.UpdateTransaction)
			r.Delete("/{id}", handlers.DeleteTransaction)
		})

		// Summary / Dashboard
		r.Get("/summary", handlers.GetSummary)
		r.Get("/summary/monthly", handlers.GetMonthlySummary)

		// Persons (legacy)
		r.HandleFunc("/persons", handlers.PersonHandler)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(common.Error(http.StatusNotFound, "Route not found"))
	})

	fmt.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
