package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sipelan/database"
	"sipelan/api"
	"sipelan/models"
	"testing"
	"time"
	"strconv"
)

func TestTransactionHandlers(t *testing.T) {
	// Setup test database and router
	database.Connect()
	api.Init()
	r := api.NewRouter()

	user := models.Person{Name: "Transaction User", Username: "txuser", Password: "password"}
	database.DB.Create(&user)
	defer database.DB.Unscoped().Delete(&user)

	category := models.Category{Name: "Test Category", PersonID: user.ID}
	database.DB.Create(&category)
	defer database.DB.Unscoped().Delete(&category)

	var transactionID uint

	t.Run("CreateTransaction_Success", func(t *testing.T) {
		input := models.Transaction{
			Date:        time.Now(),
			Description: "Test Transaksi",
			Total:       10000,
			CategoryID:  category.ID,
			Type:        "pengeluaran",
		}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(body))
		req = req.WithContext(GetTestContext(user))

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var res map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &res)
		data := res["data"].(map[string]interface{})
		transactionID = uint(data["id"].(float64))
	})

	t.Run("GetAllTransactions", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/transactions", nil)
		req = req.WithContext(GetTestContext(user))

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("GetTransactionByID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/transactions/"+strconv.Itoa(int(transactionID)), nil)
		req = req.WithContext(GetTestContext(user))

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("UpdateTransaction", func(t *testing.T) {
		input := models.Transaction{
			Date:        time.Now(),
			Description: "Updated Transaksi",
			Total:       15000,
			CategoryID:  category.ID,
			Type:        "pengeluaran",
		}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("PUT", "/api/transactions/"+strconv.Itoa(int(transactionID)), bytes.NewBuffer(body))
		req = req.WithContext(GetTestContext(user))

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("DeleteTransaction", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/transactions/"+strconv.Itoa(int(transactionID)), nil)
		req = req.WithContext(GetTestContext(user))
		
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}
