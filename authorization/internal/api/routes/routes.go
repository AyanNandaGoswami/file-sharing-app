package routes

import (
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/api/middlewares"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/api/handlers"
)

func InitializeRoutes() {
	http.HandleFunc("/authorization/v1/permission/new/", middlewares.PermissionMiddleware(handlers.RegisterNewPermission, "register_new_service"))
	http.HandleFunc("/authorization/v1/permission/all", handlers.GetAllPermission)
	http.HandleFunc("/authorization/v1/user-permission/", handlers.SetUserPermission)
}
