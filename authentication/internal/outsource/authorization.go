package outsource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/constants"
)

// Define a struct for the response to map the response JSON
type PermissionResponse struct {
	Message   string            `json:"message"`
	ExtraData map[string]string `json:"extra_data"`
}

// Define a struct to represent the request body (with primitiveUserId)
type GetUserPermissionsRequest struct {
	PrimitiveUserId string `json:"primitiveUserId"`
}

// GetUserPermissions makes a POST request with the primitiveUserId in the body
func GetUserPermissions(primitiveUserId string) (map[string]string, error) {
	// Define the API URL
	url := constants.GET_USER_PERMISSION_ENDPOINTS

	// Create the request body with primitiveUserId
	requestBody := GetUserPermissionsRequest{
		PrimitiveUserId: primitiveUserId,
	}

	// Marshal the request body to JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	// Create a new POST request with the marshaled JSON as the body
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set custom headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Unmarshal the JSON response into the PermissionResponse struct
	var permissionResponse PermissionResponse
	err = json.Unmarshal(bodyResp, &permissionResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	// Return the "extra_data" map from the response
	return permissionResponse.ExtraData, nil
}
