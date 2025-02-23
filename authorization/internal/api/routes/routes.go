package routes

import (
	"net/http"

	common_middlewares "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/middlewares"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/api/handlers"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/api/middlewares"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/constants"
)

var permissionGetter = &middlewares.PermissionGetterImplementation{}

func InitializeRoutes() {
	http.Handle(constants.ADD_NEW_ENDPOINT, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(permissionGetter)(
			http.HandlerFunc(handlers.RegisterNewAPIEndpoint),
		),
	))
	http.Handle(constants.GET_ALL_ENDPOINT, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(permissionGetter)(
			http.HandlerFunc(handlers.GetAllEndpoints),
		),
	))
	http.Handle(constants.ADD_NEW_PERMISSION, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(permissionGetter)(
			http.HandlerFunc(handlers.RegisterNewPermission),
		),
	))
	http.Handle(constants.GET_ALL_PERMISSION, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(permissionGetter)(
			http.HandlerFunc(handlers.GetAllPermission),
		),
	))
	http.Handle(constants.SET_USER_PERMISSION, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(permissionGetter)(
			http.HandlerFunc(handlers.SetUserPermission),
		),
	))

	// Internal APIs
	// Add some different middleware for validation
	http.HandleFunc(constants.GET_USER_PERMISSIONS_INTERTAL, handlers.GetUserPermissionEndpoints)
}
