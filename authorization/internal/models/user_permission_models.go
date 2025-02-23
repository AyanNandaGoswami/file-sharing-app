package models

import (
	"context"
	"fmt"
	"log"

	common_models "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/database"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userPermissionCollection *mongo.Collection

const userPermissionCollectionName = "userpermissions"

func init() {
	userPermissionCollection = database.DB.Collection(userPermissionCollectionName)
}

type UserPermission struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	User        primitive.ObjectID   `json:"user" validate:"required=This field is required."`
	Permissions []primitive.ObjectID `json:"permissions" validate:"required=This field is required."`
}

// Define a struct for the request body to hold the primitiveUserId
type GetUserPermissionsRequest struct {
	PrimitiveUserId string `json:"primitiveUserId" validate:"required=This field is required."`
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

func (userPermission *UserPermission) ValidateUserPermissionRegistrationPayload() []common_models.FieldValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(userPermission)
	var res []common_models.FieldValidationErrorResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, common_models.FieldValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	// check permissions are valid
	for _, permissionId := range userPermission.Permissions {
		validPermission := GetPermissionByID(permissionId, true)

		if validPermission == nil {
			res = append(res, common_models.FieldValidationErrorResponse{
				FieldName: "PermissionId", Message: fmt.Sprintf("Permission (%s) is inactive or invalid.", permissionId)})
		}
	}

	return res
}

func (uPermission *UserPermission) SetPermission() error {

	query := bson.M{"user": uPermission.User}
	var userPermission UserPermission

	err := userPermissionCollection.FindOne(context.Background(), query).Decode(&userPermission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no document is found, insert a new one
			userPermission = UserPermission{
				User:        uPermission.User,
				Permissions: uPermission.Permissions,
			}
			_, err := userPermissionCollection.InsertOne(context.Background(), userPermission)
			if err != nil {
				return err
			}
		} else {
			// Return any other error encountered during FindOne
			return err
		}
	} else {
		// If a document is found, update it
		update := bson.M{"permissions": uPermission.Permissions}

		_, err := userPermissionCollection.UpdateOne(context.Background(), query, bson.M{"$set": update})
		if err != nil {
			return err
		}
	}

	return nil
}

func GetUserPermissions(userId string) (map[string]string, error) {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal("Error converting string to ObjectID:", err)
	}
	query := bson.M{"user": objectId}
	var userPermission UserPermission

	err = userPermissionCollection.FindOne(context.Background(), query).Decode(&userPermission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Permission is not set for the requested user")
		} else {
			return nil, fmt.Errorf("error fetching user permissions: %v", err)
		}
	}

	// fetch the APIEndpoints for the permissions
	apiEndpoints, err := GetAPIEndpointsFromPermissions(userPermission.Permissions)
	return apiEndpoints, err

}
