package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/auth"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/models"
)

func returnErrorMessage(w http.ResponseWriter, errMessage string) {
	// Convert the error message to JSON format
	errorMessage := map[string]string{"message": errMessage}
	jsonErrorMessage, _ := json.Marshal(errorMessage)

	// Return the error message as JSON with status code 401 (Unauthorized)
	http.Error(w, string(jsonErrorMessage), http.StatusUnauthorized)
}

// AuthValidateMiddleware is a middleware to validate Headers and JWT token
func AuthValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the Authorization header is present
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			returnErrorMessage(w, "Authorization header is missing")
			return
		}

		// Split the Authorization header into "Bearer" and the token
		splitedInfo := strings.Split(authorization, " ")
		if len(splitedInfo) != 2 || splitedInfo[0] != "Bearer" {
			returnErrorMessage(w, "Invalid Authorization header format")
			return
		}

		// Extract the token from the Authorization header
		token := splitedInfo[1]

		// Validate token is blacklisted or not
		if models.IsTokenBlacklisted(token) {
			returnErrorMessage(w, "Token is not alive. Please login again.")
			return
		}

		// Retrieve user ID from JWT token
		userId, err := auth.RetrieveUserIdFromJWTToken(token)
		if err != nil {
			// Split the error message by ":"
			errorMessageParts := strings.Split(err.Error(), ":")

			// Send the error message without ":"
			returnErrorMessage(w, errorMessageParts[len(errorMessageParts)-1])
			return
		}

		// Add user id to request context if needed
		r = r.WithContext(context.WithValue(r.Context(), "userId", userId))
		r = r.WithContext(context.WithValue(r.Context(), "token", token))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
