package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sipelan/handlers"
	"sipelan/models"
	"testing"
)

func TestAuthHandlers(t *testing.T) {
	SetupTestDB()

	t.Run("Register_Success", func(t *testing.T) {
		input := models.Person{
			Username: "newuser",
			Password: "password123",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
		
		rr := httptest.NewRecorder()
		handlers.Register(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", rr.Code)
		}
	})

	t.Run("Register_Duplicate", func(t *testing.T) {
		input := models.Person{
			Username: "newuser",
			Password: "password123",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
		
		rr := httptest.NewRecorder()
		handlers.Register(rr, req)

		if rr.Code != http.StatusConflict {
			t.Errorf("Expected status 409, got %d", rr.Code)
		}
	})

	t.Run("Login_Success", func(t *testing.T) {
		// Testuser created by CreateTestUser has password "password123" and is active
		user := CreateTestUser()
		input := map[string]string{
			"username": user.Username,
			"password": "password123",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
		
		rr := httptest.NewRecorder()
		handlers.Auth(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("Login_WrongPassword", func(t *testing.T) {
		input := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
		
		rr := httptest.NewRecorder()
		handlers.Auth(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})

	t.Run("Logout", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
		rr := httptest.NewRecorder()
		handlers.Logout(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}
