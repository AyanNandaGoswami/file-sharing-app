package routes

import (
	"net/http"

	common_middlewares "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/middlewares"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/permission_manager/internal/api/handlers"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/permission_manager/internal/constants"
)

func InitializeRoutes() {
	http.Handle(constants.ADD_NEW_ENDPOINT, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.RegisterNewAPIEndpoint),
		),
	))
	http.Handle(constants.GET_ALL_ENDPOINT, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.GetAllEndpoints),
		),
	))
	http.Handle(constants.ADD_NEW_PERMISSION, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.RegisterNewPermission),
		),
	))
	http.Handle(constants.GET_ALL_PERMISSION, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.GetAllPermission),
		),
	))
	http.Handle(constants.SET_USER_PERMISSION, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.SetUserPermission),
		),
	))
}
