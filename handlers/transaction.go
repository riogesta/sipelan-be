package handlers

import (
	"encoding/json"
	"net/http"
	"sipelan/common"
	"sipelan/database"
	"sipelan/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetAllTransactions retrieves all transactions with pagination.
func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	page, limit := common.GetPaginationParams(r)
	offset := common.CalculateOffset(page, limit)

	txType := r.URL.Query().Get("type")

	person := r.Context().Value("person").(models.Person)

	var totalItems int64
	query := database.DB.Model(&models.Transaction{}).Where("person_id = ?", person.ID)
	if txType != "" {
		query = query.Where("type = ?", txType)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	var transactions []models.Transaction
	txQuery := database.DB.Preload("Category").Where("person_id = ?", person.ID).Limit(limit).Offset(offset).Order("date DESC")
	if txType != "" {
		txQuery = txQuery.Where("type = ?", txType)
	}

	if err := txQuery.Find(&transactions).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	pagination := common.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: common.CalculateTotalPages(totalItems, limit),
	}

	json.NewEncoder(w).Encode(common.NewPaginatedResponse("Success retrieving transactions", transactions, pagination))
}

// GetTransactionByID retrieves a single transaction by its ID.
func GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid ID"))
		return
	}

	person := r.Context().Value("person").(models.Person)

	var transaction models.Transaction
	if err := database.DB.Preload("Category").Where("person_id = ?", person.ID).First(&transaction, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(common.Error(http.StatusNotFound, "Transaction not found"))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Transaction found", transaction))
}

// CreateTransaction creates a new transaction.
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid request body: "+err.Error()))
		return
	}

	if transaction.Type != "pemasukan" && transaction.Type != "pengeluaran" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Type must be 'pemasukan' or 'pengeluaran'"))
		return
	}

	if transaction.Total <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Total must be greater than 0"))
		return
	}

	var category models.Category
	if err := database.DB.First(&category, transaction.CategoryID).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Category not found"))
		return
	}

	person := r.Context().Value("person").(models.Person)
	transaction.PersonID = person.ID

	result := database.DB.Create(&transaction)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, result.Error.Error()))
		return
	}

	database.DB.Preload("Category").First(&transaction, transaction.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(common.Success("Transaction created", transaction))
}

// UpdateTransaction updates an existing transaction by its ID.
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid ID"))
		return
	}

	person := r.Context().Value("person").(models.Person)

	var transaction models.Transaction
	if err := database.DB.Where("person_id = ?", person.ID).First(&transaction, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(common.Error(http.StatusNotFound, "Transaction not found"))
		return
	}

	var input models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid body"))
		return
	}

	if input.Type != "" && input.Type != "pemasukan" && input.Type != "pengeluaran" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Type must be 'pemasukan' or 'pengeluaran'"))
		return
	}

	transaction.Date = input.Date
	transaction.CategoryID = input.CategoryID
	transaction.Description = input.Description
	transaction.Total = input.Total
	transaction.Type = input.Type
	transaction.Attachment = input.Attachment

	if err := database.DB.Save(&transaction).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	database.DB.Preload("Category").First(&transaction, transaction.ID)

	json.NewEncoder(w).Encode(common.Success("Transaction updated", transaction))
}

// DeleteTransaction deletes a transaction by its ID.
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid ID"))
		return
	}

	person := r.Context().Value("person").(models.Person)

	if err := database.DB.Where("person_id = ?", person.ID).Delete(&models.Transaction{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Transaction deleted", nil))
}
