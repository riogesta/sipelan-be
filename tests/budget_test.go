package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sipelan/database"
	"sipelan/handlers"
	"sipelan/models"
	"testing"
	"time"
)

func TestBudgetHandlers(t *testing.T) {
	SetupTestDB()
	user := CreateTestUser()
	
	category := models.Category{Name: "Test Category", PersonID: user.ID}
	database.DB.Create(&category)

	t.Run("SetBudget_Create", func(t *testing.T) {
		input := map[string]interface{}{
			"category_id": category.ID,
			"amount":      500000,
			"month":       int(time.Now().Month()),
			"year":        time.Now().Year(),
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBuffer(body))
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		handlers.SetBudget(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("GetBudgetSummary", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/budgets/summary", nil)
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		handlers.GetBudgetSummary(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}
