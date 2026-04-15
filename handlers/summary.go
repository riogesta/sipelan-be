package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sipelan/common"
	"sipelan/database"
	"sipelan/models"
)

// GetSummary returns the overall income, expense, and balance totals.
func GetSummary(w http.ResponseWriter, r *http.Request) {
	var summary models.OverallSummary

	// Calculate total pemasukan
	database.DB.Model(&models.Transaction{}).
		Where("type = ?", "pemasukan").
		Select("COALESCE(SUM(total), 0)").
		Scan(&summary.TotalPemasukan)

	// Calculate total pengeluaran
	database.DB.Model(&models.Transaction{}).
		Where("type = ?", "pengeluaran").
		Select("COALESCE(SUM(total), 0)").
		Scan(&summary.TotalPengeluaran)

	summary.Balance = summary.TotalPemasukan - summary.TotalPengeluaran

	json.NewEncoder(w).Encode(common.Success("Summary retrieved", summary))
}

// GetMonthlySummary returns monthly income and expense totals for chart data.
func GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	if year == "" {
		year = "2025"
	}

	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	var results []models.MonthlySummary

	for i, monthName := range months {
		monthNum := fmt.Sprintf("%02d", i+1)

		var pemasukan, pengeluaran float64

		database.DB.Model(&models.Transaction{}).
			Where("type = ? AND strftime('%m', date) = ? AND strftime('%Y', date) = ?",
				"pemasukan", monthNum, year).
			Select("COALESCE(SUM(total), 0)").
			Scan(&pemasukan)

		database.DB.Model(&models.Transaction{}).
			Where("type = ? AND strftime('%m', date) = ? AND strftime('%Y', date) = ?",
				"pengeluaran", monthNum, year).
			Select("COALESCE(SUM(total), 0)").
			Scan(&pengeluaran)

		results = append(results, models.MonthlySummary{
			Month:       monthName,
			Pemasukan:   pemasukan,
			Pengeluaran: pengeluaran,
		})
	}

	json.NewEncoder(w).Encode(common.Success("Monthly summary retrieved", results))
}
