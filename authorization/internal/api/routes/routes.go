package routes

import (
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/api/handlers"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/constants"
)

func InitializeRoutes() {
	// Internal APIs
	// Add some different middleware for validation
	http.HandleFunc(constants.GET_USER_PERMISSIONS_INTERTAL, handlers.GetUserPermissionEndpoints)
	http.HandleFunc(constants.CHECK_AUTHORIZATION_INTERNAL, handlers.ValidateAuthorization)
}
