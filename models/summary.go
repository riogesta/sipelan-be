package models

// MonthlySummary represents the income/expense total for a specific month.
type MonthlySummary struct {
	Month       string  `json:"month"`
	Year        int     `json:"year"`
	Pemasukan   float64 `json:"pemasukan"`
	Pengeluaran float64 `json:"pengeluaran"`
}

// OverallSummary represents the total income and expense across all transactions.
type OverallSummary struct {
	TotalPemasukan   float64 `json:"total_pemasukan"`
	TotalPengeluaran float64 `json:"total_pengeluaran"`
	Balance          float64 `json:"balance"`
}
