package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sipelan/common"
	"sipelan/internal/service"
	"sipelan/models"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{s}
}

func (h *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, limit := common.GetPaginationParams(r)
	offset := common.CalculateOffset(page, limit)

	queryValues := r.URL.Query()
	filters := map[string]string{
		"type":        queryValues.Get("type"),
		"start_date":  queryValues.Get("start_date"),
		"end_date":    queryValues.Get("end_date"),
		"search":      queryValues.Get("search"),
		"category_id": queryValues.Get("category_id"),
	}

	person := r.Context().Value("person").(models.Person)

	transactions, totalItems, err := h.service.GetAllTransactions(person.ID, filters, limit, offset)
	if err != nil {
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

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}

	if validationErrors := common.ValidateStruct(transaction); validationErrors != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, common.FormatValidationError(validationErrors)))
		return
	}

	person := r.Context().Value("person").(models.Person)
	transaction.PersonID = person.ID

	if err := h.service.CreateTransaction(&transaction); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	// Preload for response
	tx, _ := h.service.GetTransactionByID(transaction.ID, person.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(common.Success("Transaction created", tx))
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	person := r.Context().Value("person").(models.Person)

	transaction, err := h.service.GetTransactionByID(uint(id), person.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(common.Error(http.StatusNotFound, "Transaction not found"))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Transaction retrieved", transaction))
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var input models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid body"))
		return
	}

	if validationErrors := common.ValidateStruct(input); validationErrors != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, common.FormatValidationError(validationErrors)))
		return
	}

	person := r.Context().Value("person").(models.Person)
	transaction, err := h.service.UpdateTransaction(uint(id), person.ID, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Transaction updated", transaction))
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	person := r.Context().Value("person").(models.Person)

	if err := h.service.DeleteTransaction(uint(id), person.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Transaction deleted", nil))
}

func (h *TransactionHandler) GetCategorySummary(w http.ResponseWriter, r *http.Request) {
	person := r.Context().Value("person").(models.Person)
	now := time.Now()
	month := fmt.Sprintf("%02d", now.Month())
	year := strconv.Itoa(now.Year())

	results, err := h.service.GetCategorySummary(person.ID, month, year)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Category summary retrieved", results))
}
