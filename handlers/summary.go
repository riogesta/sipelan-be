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

// GetSummary returns the overall income, expense, and balance totals for the authenticated user.
func GetSummary(w http.ResponseWriter, r *http.Request) {
	person := r.Context().Value("person").(models.Person)

	var totalPemasukan, totalPengeluaran float64

	database.DB.Model(&models.Transaction{}).Where("person_id = ? AND type = ?", person.ID, "pemasukan").Select("COALESCE(SUM(total), 0)").Scan(&totalPemasukan)
	database.DB.Model(&models.Transaction{}).Where("person_id = ? AND type = ?", person.ID, "pengeluaran").Select("COALESCE(SUM(total), 0)").Scan(&totalPengeluaran)

	summary := models.OverallSummary{
		TotalPemasukan:   totalPemasukan,
		TotalPengeluaran: totalPengeluaran,
		Balance:          totalPemasukan - totalPengeluaran,
	}

	json.NewEncoder(w).Encode(common.Success("Summary retrieved", summary))
}

// GetMonthlySummary returns monthly income and expense totals for the authenticated user.
func GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	person := r.Context().Value("person").(models.Person)
	
	year := r.URL.Query().Get("year")
	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}

	months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	var results []models.MonthlySummary

	for _, m := range months {
		var pemasukan, pengeluaran float64
		database.DB.Model(&models.Transaction{}).Where("person_id = ? AND type = ? AND strftime('%m', date) = ? AND strftime('%Y', date) = ?", person.ID, "pemasukan", m, year).Select("COALESCE(SUM(total), 0)").Scan(&pemasukan)
		database.DB.Model(&models.Transaction{}).Where("person_id = ? AND type = ? AND strftime('%m', date) = ? AND strftime('%Y', date) = ?", person.ID, "pengeluaran", m, year).Select("COALESCE(SUM(total), 0)").Scan(&pengeluaran)

		monthInt, _ := strconv.Atoi(m)
		monthName := time.Month(monthInt).String()

		results = append(results, models.MonthlySummary{
			Month:       monthName,
			Pemasukan:   pemasukan,
			Pengeluaran: pengeluaran,
		})
	}

	json.NewEncoder(w).Encode(common.Success("Monthly summary retrieved", results))
}

// GetChartData returns flexible chart data based on the 'view' parameter (daily, weekly, monthly) for the authenticated user.
func GetChartData(w http.ResponseWriter, r *http.Request) {
	person := r.Context().Value("person").(models.Person)
	
	view := r.URL.Query().Get("view") // daily, weekly, monthly
	if view == "" {
		view = "monthly"
	}

	var results []models.ChartData

	if view == "monthly" {
		year := r.URL.Query().Get("year")
		if year == "" {
			year = strconv.Itoa(time.Now().Year())
		}
		months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		for i, monthName := range months {
			monthNum := fmt.Sprintf("%02d", i+1)
			var in, out float64
			database.DB.Model(&models.Transaction{}).
				Where("person_id = ? AND type = ? AND strftime('%m', date) = ? AND strftime('%Y', date) = ?", person.ID, "pemasukan", monthNum, year).
				Select("COALESCE(SUM(total), 0)").Scan(&in)
			database.DB.Model(&models.Transaction{}).
				Where("person_id = ? AND type = ? AND strftime('%m', date) = ? AND strftime('%Y', date) = ?", person.ID, "pengeluaran", monthNum, year).
				Select("COALESCE(SUM(total), 0)").Scan(&out)
			results = append(results, models.ChartData{Label: monthName, Pemasukan: in, Pengeluaran: out})
		}
	} else if view == "weekly" {
		// Get last 7 days including today
		for i := 6; i >= 0; i-- {
			var in, out float64
			var label string
			
			// Get date string for the query (e.g., '2024-04-26')
			queryDate := ""
			database.DB.Raw("SELECT date('now', 'localtime', ?)", fmt.Sprintf("-%d day", i)).Scan(&queryDate)
			// Get formatted label
			database.DB.Raw("SELECT strftime('%d/%m', date('now', 'localtime', ?))", fmt.Sprintf("-%d day", i)).Scan(&label)

			database.DB.Model(&models.Transaction{}).
				Where("person_id = ? AND type = ? AND date(date) = ?", person.ID, "pemasukan", queryDate).
				Select("COALESCE(SUM(total), 0)").Scan(&in)
			database.DB.Model(&models.Transaction{}).
				Where("person_id = ? AND type = ? AND date(date) = ?", person.ID, "pengeluaran", queryDate).
				Select("COALESCE(SUM(total), 0)").Scan(&out)
			
			results = append(results, models.ChartData{Label: label, Pemasukan: in, Pengeluaran: out})
		}
	} else if view == "daily" {
		// Group by day for the CURRENT month FOR THIS PERSON
		rows, err := database.DB.Raw(`
			SELECT 
				strftime('%d', date) as day_label, 
				SUM(CASE WHEN type = 'pemasukan' THEN total ELSE 0 END) as in_total,
				SUM(CASE WHEN type = 'pengeluaran' THEN total ELSE 0 END) as out_total
			FROM transactions 
			WHERE person_id = ? 
			  AND strftime('%m', date) = strftime('%m', 'now', 'localtime') 
			  AND strftime('%Y', date) = strftime('%Y', 'now', 'localtime')
			GROUP BY day_label 
			ORDER BY day_label ASC
		`, person.ID).Rows()
		
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var label string
				var in, out float64
				rows.Scan(&label, &in, &out)
				results = append(results, models.ChartData{Label: label, Pemasukan: in, Pengeluaran: out})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(common.Success("Chart data retrieved", results))
}
