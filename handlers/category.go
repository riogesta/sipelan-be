package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sipelan/common"
	"sipelan/database"
	"sipelan/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetAllCategories retrieves all categories with pagination.
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	page, limit := common.GetPaginationParams(r)
	offset := common.CalculateOffset(page, limit)

	person := r.Context().Value("person").(models.Person)

	var totalItems int64
	if err := database.DB.Model(&models.Category{}).Where("person_id = ?", person.ID).Count(&totalItems).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	var categories []models.Category
	if err := database.DB.Where("person_id = ?", person.ID).Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
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

	json.NewEncoder(w).Encode(common.NewPaginatedResponse("Success retrieving categories", categories, pagination))
}

// GetCategoryByID retrieves a single category by its ID.
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid ID"))
		return
	}

	person := r.Context().Value("person").(models.Person)

	var category models.Category
	if err := database.DB.Where("person_id = ?", person.ID).First(&category, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(common.Error(http.StatusNotFound, "Category not found"))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Category found", category))
}

// CreateCategory creates a new category.
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}

	if category.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Name is required"))
		return
	}

	person := r.Context().Value("person").(models.Person)
	category.PersonID = person.ID

	// validation: check if category with the same name already exists for this person (case-insensitive)
	var existingCategory models.Category
	if err := database.DB.Where("person_id = ? AND LOWER(name) = LOWER(?)", person.ID, category.Name).First(&existingCategory).Error; err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(common.Error(http.StatusConflict, "Category with this name already exists"))
		return
	}

	if result := database.DB.Create(&category); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, result.Error.Error()))
		return
	}

	slog.Info("category created")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(common.Success("Category created", category))
}

// UpdateCategory updates an existing category by its ID.
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid ID"))
		return
	}

	person := r.Context().Value("person").(models.Person)

	var category models.Category
	if err := database.DB.Where("person_id = ?", person.ID).First(&category, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(common.Error(http.StatusNotFound, "Category not found"))
		return
	}

	var input models.Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid body"))
		return
	}

	category.Name = input.Name
	category.Description = input.Description

	if err := database.DB.Save(&category).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Category updated", category))
}

// DeleteCategory deletes a category by its ID.
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid ID"))
		return
	}

	person := r.Context().Value("person").(models.Person)

	if err := database.DB.Where("person_id = ?", person.ID).Delete(&models.Category{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	json.NewEncoder(w).Encode(common.Success("Category deleted", nil))
}
