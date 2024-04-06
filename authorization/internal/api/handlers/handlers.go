package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/models"
)

func RegisterNewService(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	var service models.Service

	json.NewDecoder(r.Body).Decode(&service)

	validation_err := service.ValidateServiceRegistrationPayload()
	if validation_err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validation_err)
		return
	}

	err := service.RegisterService()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.APIResponse{Message: "Successfully registered new service."})

}
