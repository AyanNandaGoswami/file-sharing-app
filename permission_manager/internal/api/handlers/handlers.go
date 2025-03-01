package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	common_models "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/permission_manager/internal/models"
)

func RegisterNewAPIEndpoint(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	var api_endpoints models.APIEndpoints
	json.NewDecoder(r.Body).Decode(&api_endpoints)

	if validatonErr := api_endpoints.ValidateAPIEndpointsRegistrationPayload(); validatonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validatonErr)
		return
	}

	if err := api_endpoints.RegisterNewAPIEndpoint(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common_models.APIResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(common_models.APIResponse{Message: "Successfully added the endpoint."})

}

func GetAllEndpoints(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	endpoinds, err := models.AllAPIEndpoints()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common_models.APIResponse{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(endpoinds)
}

func RegisterNewPermission(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	var permission models.Permission

	json.NewDecoder(r.Body).Decode(&permission)

	if validationErr := permission.ValidatePermissionRegistrationPayload(); validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErr)
		return
	}

	err := permission.CreatePermission()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common_models.APIResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(common_models.APIResponse{Message: "Successfully added new permission."})

}

func GetAllPermission(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	// Extract the isActive from the query parameters
	isActive := r.URL.Query().Get("isActive")
	var isActiveAck *bool

	if isActive != "" {
		pattern := `^[0-1]$`
		matched, err := regexp.MatchString(pattern, isActive)
		if err != nil || !matched {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common_models.APIResponse{Message: "isActive value is invalid."})
			return
		}

		isActiveBool := false

		if isActive == "1" {
			isActiveBool = true
		}
		isActiveAck = &isActiveBool
	}

	permissions, err := models.AllPermissions(isActiveAck)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common_models.APIResponse{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permissions)
}

func SetUserPermission(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	var userPermission models.UserPermission

	json.NewDecoder(r.Body).Decode(&userPermission)

	if validationErrors := userPermission.ValidateUserPermissionRegistrationPayload(); validationErrors != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrors)
		return
	}

	err := userPermission.SetPermission()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common_models.APIResponse{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common_models.APIResponse{Message: "User permission has been updated successfully."})
}
