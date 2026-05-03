package api

import (
	"sipelan/handlers"
	"sipelan/internal/repository"
	"sipelan/internal/service"
	"sipelan/database"
	"sipelan/common"
	"sipelan/config"

	"github.com/go-chi/chi/v5"
)

var (
	TxRepo repository.TransactionRepository
	TxService service.TransactionService
	TxHandler *handlers.TransactionHandler
)

func Init() {
	TxRepo = repository.NewTransactionRepository(database.DB)
	TxService = service.NewTransactionService(TxRepo)
	TxHandler = handlers.NewTransactionHandler(TxService)
}

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware MUST be defined before any routes
	r.Use(common.CORSMiddleware(cfg.AllowedOrigins))
	r.Use(common.JSONMiddleware)

	// Initialize layers if not already (for tests)
	if TxHandler == nil {
		Init()
	}

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handlers.Auth)
		r.Post("/register", handlers.Register)
		r.Post("/logout", handlers.Logout)
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
			r.Get("/", TxHandler.GetAll)
			r.Post("/", TxHandler.Create)
			r.Get("/{id}", TxHandler.GetByID)
			r.Put("/{id}", TxHandler.Update)
			r.Delete("/{id}", TxHandler.Delete)
		})

		// Summary / Dashboard
		r.Get("/summary", handlers.GetSummary)
		r.Get("/summary/monthly", handlers.GetMonthlySummary)
		r.Get("/summary/chart", handlers.GetChartData)
		r.Get("/summary/categories", TxHandler.GetCategorySummary)
		r.Get("/summary/budget", handlers.GetBudgetSummary)

		// Budgets
		r.Post("/budgets", handlers.SetBudget)

		// Upload
		r.Post("/upload", handlers.UploadFile)

		// Persons (legacy)
		r.HandleFunc("/persons", handlers.PersonHandler)
	})

	return r
}
