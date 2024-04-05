package models

type FielValidationErrorResponse struct {
	FieldName string `json:"field_name"`
	Message   string `json:"message"`
}

type APIResponse struct {
	Message   string      `json:"message"`
	ExtraData interface{} `json:"extra_data"`
}
