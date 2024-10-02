package routes

import (
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/nexus/internal/api/handlers"
)

func InitializeRoutes() {
	http.HandleFunc("/nexus/v1/service/register/", handlers.RegisterNewService)
	http.HandleFunc("/nexus/v1/service/all/", handlers.GetAllServices)
	http.HandleFunc("/nexus/v1/service/", handlers.ValidateServiceId)
}
