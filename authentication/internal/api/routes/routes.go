package routes

import (
	"net/http"

	common_middlewares "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/middlewares"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/api/handlers"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/constants"
)

func InitializeRoutes() {
	http.HandleFunc(constants.REGISTER, handlers.RegisterNewUser)
	http.HandleFunc(constants.LOGIN, handlers.Login)

	// Require authentication for the following endpoints
	http.Handle(constants.GET_USER_INFO, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.UserDetail),
		),
	))
	http.Handle(constants.UPDATE_USER_INFO, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.UpdateUserInfo),
		),
	))
	http.Handle(constants.DELETE_USER_ACCOUNT, common_middlewares.AuthValidateMiddleware(
		common_middlewares.PermissionValidationMiddleware(
			http.HandlerFunc(handlers.DeleteUser),
		),
	))
	http.Handle(constants.LOGOUT, common_middlewares.AuthValidateMiddleware(http.HandlerFunc(handlers.Logout)))
}
