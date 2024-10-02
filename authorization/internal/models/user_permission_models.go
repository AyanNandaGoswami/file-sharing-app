package models

import (
	"context"
	"fmt"

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
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	User            primitive.ObjectID `json:"user" validate:"required=This field is required."`
	PermissionNames []string           `json:"permissionnames" validate:"required=This field is required."`
}

type PermissionWithAction struct {
	Action       string             `json:"action" validate:"required=This field is required."`
	PermissionId primitive.ObjectID `json:"permissionid" validate:"required=This field is required."`
}

type CreateUserPermission struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty"`
	User            primitive.ObjectID     `json:"user" validate:"required=This field is required."`
	PermissionData  []PermissionWithAction `json:"permissionwithaction" validate:"required=This field is required."`
	PermissionNames []string               `json:"permissionnames"`
}

func (up *CreateUserPermission) ValidateUserPermissionRegistrationPayload() []FieldValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(up)
	var res []FieldValidationErrorResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, FieldValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	// check permissions are valid
	for _, permissionItem := range up.PermissionData {
		permissionObj, validPermission := GetPermissionByID(permissionItem.PermissionId, true)

		if !validPermission {
			res = append(res, FieldValidationErrorResponse{
				FieldName: "PermissionId", Message: fmt.Sprintf("Permission (%s) is inactive or invalid.", permissionItem.PermissionId)})
		} else {
			up.PermissionNames = append(up.PermissionNames, permissionObj.Name)
		}
	}

	return res
}

func (up *CreateUserPermission) SetPermission() error {

	query := bson.M{"user": up.User}
	var userPermission UserPermission

	err := userPermissionCollection.FindOne(context.Background(), query).Decode(&userPermission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no document is found, insert a new one
			userPermission = UserPermission{
				User:            up.User,
				PermissionNames: up.PermissionNames,
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
		update := bson.M{"PermissionNames": up.PermissionNames}

		_, err := userPermissionCollection.UpdateOne(context.Background(), query, bson.M{"$set": update})
		if err != nil {
			return err
		}
	}

	return nil
}
