package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sipelan/common"
	"sipelan/database"
	"sipelan/models"
)

func PersonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		createPerson(w, r)
		return
	}

	getPeople(w, r)
}

func createPerson(w http.ResponseWriter, r *http.Request) {

	var person models.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}

	if person.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Name is required"))
		return
	}

	if person.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Error(http.StatusBadRequest, "Password is required"))
		return
	}

	hashedPassword, err := common.HashPassword(person.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, "Error hashing password"))
		return
	}

	person.Password = hashedPassword

	result := database.DB.Create(&person)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Error(http.StatusInternalServerError, "Error saving to database"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common.Success(fmt.Sprintf("Hello %s, you have been saved to the database!", person.Name), person))
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []models.Person
	database.DB.Find(&people)

	json.NewEncoder(w).Encode(common.Success("List of People", people))
}
