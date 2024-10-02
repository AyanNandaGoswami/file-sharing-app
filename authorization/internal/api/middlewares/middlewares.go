package middlewares

import "net/http"

func PermissionMiddleware(next http.Handler, requiredPermissionCode string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userPermissions := getUserPermissions(r)

		if !hasPermission(userPermissions, requiredPermissionCode) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getUserPermissions(userId string) []string {
	// Implement logic to retrieve user permissions from the request context/session
	// Example hardcoded permissions for demonstration:
	return []string{"register_new_service", "view_all_permissions", "set_user_permission"}
}

func hasPermission(userPermissions []string, requiredPermissionCode string) bool {
	for _, perm := range userPermissions {
		if perm == requiredPermissionCode {
			return true
		}
	}
	return false
}
