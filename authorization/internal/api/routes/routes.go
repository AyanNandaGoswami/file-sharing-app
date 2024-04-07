package routes

import (
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/api/handlers"
)

func InitializeRoutes() {
	http.HandleFunc("/authorization/v1/service/register/", handlers.RegisterNewService)
	http.HandleFunc("/authorization/v1/permission/new/", handlers.RegisterNewPermission)
}
