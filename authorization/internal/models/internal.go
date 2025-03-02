package models

import (
	common_models "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/models"
	"github.com/go-playground/validator/v10"
)

// Define a struct for the request body to hold the primitiveUserId
type GetUserPermissionsRequest struct {
	PrimitiveUserId string `json:"primitiveUserId" validate:"required=This field is required."`
}

type PermissionValidation struct {
	PrimitiveUserId string `json:"primitiveUserId"`
	Token           string `json:"token"`
	ValidateBy      string `json:"validate_by" validate:"required=This field is required,oneof=token primitiveUserId"`
	RequestedUrl    string `json:"requestedUrl" validate:"required=This field is required."`
	RequestedMethod string `json:"requestedMethod" validate:"required=This field is required."`
	ReturnUserId    bool   `json:"returnUserId"`
}

func (up *GetUserPermissionsRequest) ValidateGetUserPermissionsRequestRegistrationPayload() []common_models.FieldValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(up)
	var res []common_models.FieldValidationErrorResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, common_models.FieldValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	return res
}

func (pv *PermissionValidation) ValidatePermissionValidadtionPayload() []common_models.FieldValidationErrorResponse {
	validate := validator.New()

	// Perform conditional validation based on ValidateBy field
	if pv.ValidateBy == "token" {
		// Validate Token is required when ValidateBy is "token"
		if pv.Token == "" {
			return []common_models.FieldValidationErrorResponse{
				{
					FieldName: "token",
					Message:   "This field is required when validate_by is token.",
				},
			}
		}
	} else if pv.ValidateBy == "primitiveUserId" {
		// Validate PrimitiveUserId is required when ValidateBy is "user_id"
		if pv.PrimitiveUserId == "" {
			return []common_models.FieldValidationErrorResponse{
				{
					FieldName: "primitiveUserId",
					Message:   "This field is required when validate_by is primitiveUserId.",
				},
			}
		}
	}

	// General validation for other fields
	err := validate.Struct(pv)
	var res []common_models.FieldValidationErrorResponse

	if err != nil {
		// Loop through the validation errors
		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, common_models.FieldValidationErrorResponse{
				FieldName: err.StructField(),
				Message:   err.Param(),
			})
		}
	}

	return res
}
