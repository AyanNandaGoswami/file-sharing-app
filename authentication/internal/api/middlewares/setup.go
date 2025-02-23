package middlewares

import "github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/outsource"

type PermissionGetterImplementation struct{}

func (p *PermissionGetterImplementation) GetUserPermissionEndpoints(token string) (map[string]string, error) {
	// Define the GetUserPermissionEndpoints method
	apiEndpoints, err := outsource.GetUserPermissions(token)
	if err != nil {
		return nil, err
	}
	return apiEndpoints, nil
}
