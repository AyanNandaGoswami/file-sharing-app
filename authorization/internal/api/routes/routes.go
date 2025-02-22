package routes

import (
	"net/http"

	common_middlewares "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/middlewares"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/api/handlers"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/api/middlewares"
)

func InitializeRoutes() {
	http.Handle("/authorization/v1/endpoint/add/", common_middlewares.AuthValidateMiddleware(
		middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.RegisterNewAPIEndpoint),
		),
	))
	http.Handle("/authorization/v1/endpoint/all", common_middlewares.AuthValidateMiddleware(
		middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.GetAllEndpoints),
		),
	))
	http.Handle("/authorization/v1/permission/add/", common_middlewares.AuthValidateMiddleware(
		middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.RegisterNewPermission),
		),
	))
	http.Handle("/authorization/v1/permission/all", common_middlewares.AuthValidateMiddleware(
		middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.GetAllPermission),
		),
	))
	http.Handle("/authorization/v1/user-permission/set/", common_middlewares.AuthValidateMiddleware(
		middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.SetUserPermission),
		),
	))
}
