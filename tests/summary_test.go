package tests

import (
	"net/http"
	"net/http/httptest"
	"sipelan/handlers"
	"testing"
)

func TestSummaryHandlers(t *testing.T) {
	SetupTestDB()
	user := CreateTestUser()

	t.Run("GetSummary", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		handlers.GetSummary(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("GetMonthlySummary", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/summary/monthly", nil)
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		handlers.GetMonthlySummary(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("GetChartData", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/summary/chart?view=monthly", nil)
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		handlers.GetChartData(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}
