package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/auth"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/models"
)

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	// Decode the user payload from the request body
	var user models.CreateUser
	json.NewDecoder(r.Body).Decode(&user)

	// Validate user registration payload
	err := user.ValidateUserRegistrationPayload()

	// If there is an error, return 400 bad request
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Create new user
	if user.CreateNewUser() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.FielValidationErrorResponse{Message: "Account does not created, contact to website admin."})
		return
	}

	// Return successful response
	json.NewEncoder(w).Encode(models.APIResponse{Message: "Account created successfully.", ExtraData: nil})
}

func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	// Check if request method is PUT
	if r.Method != "PUT" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	// Decode the user payload from the request body
	var user models.UserUpdate
	json.NewDecoder(r.Body).Decode(&user)

	// Validate user registration payload
	validation_err := user.ValidateUserUpdatePayload()

	if validation_err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validation_err)
		return
	}

	userId := r.Context().Value("userId").(string)

	if err := user.UpdateUserByID(userId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{Message: err.Error()})
		return
	}

	// Return successful response
	json.NewEncoder(w).Encode(models.APIResponse{Message: "User details updated successfully.", ExtraData: nil})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Check if request method is DELETE
	if r.Method != "DELETE" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userId := r.Context().Value("userId").(string)

	err := models.DeleteUser(userId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{Message: err.Error(), ExtraData: nil})
		return
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(models.APIResponse{Message: "Successfully deleted the user."})
}

func Login(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	var login models.Login
	json.NewDecoder(r.Body).Decode(&login)

	if err_payload := login.ValiadteLoginPayload(); err_payload != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err_payload)
		return
	}

	if userId, err := login.Authenticate(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.APIResponse{Message: "Invalid authentication credentials."})
		return
	} else {
		token, err := auth.GenerateNewJWToken(userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.APIResponse{Message: "Something went wrong!"})
		} else {
			json.NewEncoder(w).Encode(models.APIResponse{Message: "Logged in successfully.",
				ExtraData: map[string]string{"access_token": token}})
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Set the content type to json
	w.Header().Set("Content-Type", "application/json")

	token := r.Context().Value("token").(string)

	models.NewBlacklistedToken(token)

	w.WriteHeader(http.StatusOK)
}

func UserDetail(w http.ResponseWriter, r *http.Request) {
	// Check if request method is GET
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value("userId").(string)

	user, err := models.GetUserByUUID(userId)
	if err != nil {
		log.Fatal(err)
	} else {
		json.NewEncoder(w).Encode(user)
	}
}
