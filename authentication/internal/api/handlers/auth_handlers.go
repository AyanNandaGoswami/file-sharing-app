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

// func ValidateAccesssToken(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	_, claims, _ := jwtauth.FromContext(r.Context())
// 	userId, _ := claims["user_id"].(string)

// 	json.NewEncoder(w).Encode(map[string]string{"user_id": userId})

// }
