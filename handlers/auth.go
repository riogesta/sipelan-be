package handlers

import (
	"encoding/json"
	"net/http"
	"sipelan/common"
	"sipelan/database"
	"sipelan/models"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	var req models.Person

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}

	var person models.Person
	result := database.DB.Where("username = ?", req.Username).First(&person)

	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.Error(http.StatusUnauthorized, "Invalid username or password"))
		return
	}

	if !common.VerifyPassword(req.Password, person.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.Error(http.StatusUnauthorized, "Invalid username or password"))
		return
	}

	// Check if user is active
	if !person.IsActive {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(common.Error(http.StatusForbidden, "Your account is not active. Please contact administrator."))
		return
	}

	token, err := common.CreateToken(person.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, "Failed to create token"))
		return
	}

	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to false for local development over HTTP
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common.Success("Authentication successful", map[string]interface{}{
		"token":     token,
		"person_id": person.ID,
		"username":  person.Username,
	}))
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.Person

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}

	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Username and password are required"))
		return
	}

	// Check if user already exists
	var existing models.Person
	if result := database.DB.Where("username = ?", req.Username).First(&existing); result.Error == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(common.Error(http.StatusConflict, "Username already exists"))
		return
	}

	// Hash password
	hashedPassword, err := common.HashPassword(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, "Failed to process password"))
		return
	}
	req.Password = hashedPassword
	req.IsActive = false // Ensure explicit false on registration

	if result := database.DB.Create(&req); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, result.Error.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(common.Success("Account created successfully. Please wait for activation.", nil))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // Expire immediately
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common.Success("Logout successful", nil))
}
