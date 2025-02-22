package routes

import (
	"net/http"

	middlewares "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/middlewares"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/api/handlers"
)

func InitializeRoutes() {
	http.HandleFunc("/auth/v1/register/", handlers.RegisterNewUser)
	http.HandleFunc("/auth/v1/login/", handlers.Login)

	// Require authentication for the following endpoints
	http.Handle("/auth/v1/userinfo", middlewares.AuthValidateMiddleware(http.HandlerFunc(handlers.UserDetail)))
	http.Handle("/auth/v1/logout/", middlewares.AuthValidateMiddleware(http.HandlerFunc(handlers.Logout)))
	http.Handle("/auth/v1/update/userinfo/", middlewares.AuthValidateMiddleware(http.HandlerFunc(handlers.UpdateUserInfo)))
	http.Handle("/auth/v1/account/delete/", middlewares.AuthValidateMiddleware(http.HandlerFunc(handlers.DeleteUser)))
}
