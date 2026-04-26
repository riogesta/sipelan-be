package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sipelan/common"
	"sipelan/database"
	"sipelan/models"
	"strconv"
	"time"
)

// SetBudget creates or updates a budget for a category
func SetBudget(w http.ResponseWriter, r *http.Request) {
	person := r.Context().Value("person").(models.Person)
	
	var req struct {
		CategoryID uint    `json:"category_id"`
		Amount     float64 `json:"amount"`
		Month      int     `json:"month"`
		Year       int     `json:"year"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}

	// Validate category ownership
	var category models.Category
	if err := database.DB.Where("id = ? AND person_id = ?", req.CategoryID, person.ID).First(&category).Error; err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(common.Error(http.StatusForbidden, "Category not found or access denied"))
		return
	}

	// Default to current month/year if not provided
	now := time.Now()
	if req.Month == 0 {
		req.Month = int(now.Month())
	}
	if req.Year == 0 {
		req.Year = now.Year()
	}

	var budget models.Budget
	result := database.DB.Where("category_id = ? AND month = ? AND year = ? AND person_id = ?", req.CategoryID, req.Month, req.Year, person.ID).First(&budget)

	if result.Error == nil {
		// Update existing
		budget.Amount = req.Amount
		database.DB.Save(&budget)
	} else {
		// Create new
		budget = models.Budget{
			PersonID:   person.ID,
			CategoryID: req.CategoryID,
			Amount:     req.Amount,
			Month:      req.Month,
			Year:       req.Year,
		}
		database.DB.Create(&budget)
	}

	json.NewEncoder(w).Encode(common.Success("Budget set successfully", budget))
}

// GetBudgetSummary returns budgets with their current spending
func GetBudgetSummary(w http.ResponseWriter, r *http.Request) {
	person := r.Context().Value("person").(models.Person)
	
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	now := time.Now()
	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	if month == 0 {
		month = int(now.Month())
	}
	if year == 0 {
		year = now.Year()
	}

	var budgets []models.Budget
	database.DB.Preload("Category").Where("month = ? AND year = ? AND person_id = ?", month, year, person.ID).Find(&budgets)

	type BudgetUsage struct {
		models.Budget
		Used       float64 `json:"used"`
		Percentage float64 `json:"percentage"`
	}

	var results []BudgetUsage
	for _, b := range budgets {
		var used float64
		// Sum expenses for this category in the given month/year FOR THIS PERSON
		database.DB.Model(&models.Transaction{}).
			Where("person_id = ? AND category_id = ? AND type = 'pengeluaran' AND strftime('%m', date) = ? AND strftime('%Y', date) = ?",
				person.ID, b.CategoryID, fmt.Sprintf("%02d", month), strconv.Itoa(year)).
			Select("COALESCE(SUM(total), 0)").
			Scan(&used)

		percentage := 0.0
		if b.Amount > 0 {
			percentage = (used / b.Amount) * 100
		}

		results = append(results, BudgetUsage{
			Budget:     b,
			Used:       used,
			Percentage: percentage,
		})
	}

	json.NewEncoder(w).Encode(common.Success("Budget summary retrieved", results))
}
