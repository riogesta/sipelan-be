package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sipelan/handlers"
	"sipelan/models"
	"testing"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func TestCategoryHandlers(t *testing.T) {
	SetupTestDB()
	user := CreateTestUser()

	var categoryID uint

	t.Run("CreateCategory_Success", func(t *testing.T) {
		input := models.Category{
			Name: "Makanan",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		handlers.CreateCategory(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", rr.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		categoryID = uint(data["id"].(float64))
	})

	t.Run("GetAllCategories", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		handlers.GetAllCategories(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("UpdateCategory", func(t *testing.T) {
		input := models.Category{
			Name: "Makanan & Minuman",
		}
		body, _ := json.Marshal(input)
		
		// Setup router to handle URL param
		r := chi.NewRouter()
		r.Put("/api/categories/{id}", handlers.UpdateCategory)
		
		req := httptest.NewRequest(http.MethodPut, "/api/categories/"+strconv.Itoa(int(categoryID)), bytes.NewBuffer(body))
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("DeleteCategory", func(t *testing.T) {
		r := chi.NewRouter()
		r.Delete("/api/categories/{id}", handlers.DeleteCategory)
		
		req := httptest.NewRequest(http.MethodDelete, "/api/categories/"+strconv.Itoa(int(categoryID)), nil)
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}
