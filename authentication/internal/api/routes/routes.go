package routes

import (
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/api/handlers"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/api/middlewares"
)

func InitializeRoutes() {
	http.HandleFunc("/auth/v1/register/", handlers.RegisterNewUser)
	http.HandleFunc("/auth/v1/login/", handlers.Login)

	http.Handle("/auth/v1/user-detail", middlewares.AuthValidateMiddleware(http.HandlerFunc(handlers.UserDetail)))
}
