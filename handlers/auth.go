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

	token, err := common.CreateToken(person.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, "Failed to create token"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common.Success("Authentication successful", map[string]string{
		"token": token,
	}))
}
