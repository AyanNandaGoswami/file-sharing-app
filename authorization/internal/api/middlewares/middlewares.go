package middlewares

import (
	"encoding/json"
	"net/http"

	common_middlewares "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/middlewares"
	common_models "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/models"
)

func ReturnErrorMessage(w http.ResponseWriter, errMessage string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(common_models.APIResponse{Message: errMessage, ExtraData: nil})
}

func PermissionValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		primitiveUserId := r.Context().Value(common_middlewares.PrimitiveUserIdKey).(string)
		permissionEndpoints, err := getUserPermissionEndpoints(primitiveUserId)
		if err != nil {
			ReturnErrorMessage(w, err.Error(), 400)
			return
		}

		if !hasPermission(permissionEndpoints, r.URL.Path, r.Method) {
			ReturnErrorMessage(w, "You do not have permission to perform this action", 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getUserPermissionEndpoints(userId string) (map[string]string, error) {
	apiEndpoints, err := models.GetUserPermissions(userId)
	if err != nil {
		return nil, err
	}
	return apiEndpoints, nil
}

func hasPermission(userPermissionEndpoints map[string]string, requestedUrl string, requesedMethod string) bool {
	method, exists := userPermissionEndpoints[requestedUrl]
	if exists {
		return method == requesedMethod
	}
	return false
}
