package constants

import "fmt"

const AUTHORIZATION_SERVICE_BASE_URL = "http://127.0.0.1:4002"

var GET_USER_PERMISSION_ENDPOINTS = fmt.Sprintf("%s/authorization/v1/user-permission/get/", AUTHORIZATION_SERVICE_BASE_URL)
