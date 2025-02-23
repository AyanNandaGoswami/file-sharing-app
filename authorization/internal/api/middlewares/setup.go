package middlewares

import "github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/models"

type PermissionGetterImplementation struct{}

func (p *PermissionGetterImplementation) GetUserPermissionEndpoints(primitiveUserId string) (map[string]string, error) {
	// Define the GetUserPermissionEndpoints method
	apiEndpoints, err := models.GetUserPermissions(primitiveUserId)
	if err != nil {
		return nil, err
	}
	return apiEndpoints, nil
}
