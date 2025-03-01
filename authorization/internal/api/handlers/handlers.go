package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	common_models "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/models"
	auth "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/utilities"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/models"
)

func GetUserPermissionEndpoints(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	var requestData models.GetUserPermissionsRequest
	json.NewDecoder(r.Body).Decode(&requestData)
	if validationErrors := requestData.ValidateGetUserPermissionsRequestRegistrationPayload(); validationErrors != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrors)
		return
	}

	// Call the method to get user permissions
	apiEndpoints, err := models.GetUserPermissions(requestData.PrimitiveUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common_models.APIResponse{Message: err.Error(), ExtraData: nil})
		return
	}

	// Send the response with the permission data
	json.NewEncoder(w).Encode(common_models.APIResponse{Message: "Ok", ExtraData: apiEndpoints})
}

func ValidateAuthorization(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	var permissionValidationPayload models.PermissionValidation
	var primitiveUserId string
	json.NewDecoder(r.Body).Decode(&permissionValidationPayload)

	if validationErrors := permissionValidationPayload.ValidatePermissionValidadtionPayload(); validationErrors != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrors)
		return
	}

	// if validate_by token then fetch the PrimitiveUserId from the token
	if permissionValidationPayload.ValidateBy == "token" {
		info, err := auth.RetrieveDetilsFromJWT(permissionValidationPayload.Token)
		if err != nil {
			// Split the error message by ":"
			errorMessageParts := strings.Split(err.Error(), ":")

			w.WriteHeader(http.StatusUnauthorized)
			// Send the error message without ":"
			json.NewEncoder(w).Encode(common_models.APIResponse{
				Message:   errorMessageParts[len(errorMessageParts)-1],
				ExtraData: map[string]bool{"authorized": false}})
			return
		}
		primitiveUserId = info.PrimitiveUserId
	} else {
		primitiveUserId = permissionValidationPayload.PrimitiveUserId
	}

	// Call the method to get user permissions
	apiEndpoints, err := models.GetUserPermissions(primitiveUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common_models.APIResponse{Message: err.Error(), ExtraData: nil})
		return
	}

	// check requested url is exists or not in the apiEndpoints with the same method
	method, exists := apiEndpoints[permissionValidationPayload.RequestedUrl]
	if exists {
		if method == permissionValidationPayload.RequestedMethod {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(common_models.APIResponse{Message: "Success", ExtraData: map[string]bool{"authorized": true}})
			return
		}
	}
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(common_models.APIResponse{
		Message:   "You don't have permission to perform this action.",
		ExtraData: map[string]bool{"authorized": false}})

}
