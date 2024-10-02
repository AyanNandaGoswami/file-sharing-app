package models

import (
	"context"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/database"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var permissionCollection *mongo.Collection

const permissionCollectionName = "permissions"

func init() {
	permissionCollection = database.DB.Collection(permissionCollectionName)
}

type Permission struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" validate:"required=This field is required."`
	Service  primitive.ObjectID `bson:"service" validate:"required=This field is required."`
	IsActive bool               `json:"is_active" validate:"required=This field is required."`
}

type PermissionList struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" validate:"required=This field is required."`
	IsActive bool               `json:"is_active" validate:"required=This field is required."`
}

func (p *Permission) ValidatePermissionRegistrationPayload() []FieldValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(p)
	var res []FieldValidationErrorResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, FieldValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	// validate permission is already exists or not
	var result Permission
	query := bson.M{"name": p.Name}

	if err := permissionCollection.FindOne(context.Background(), query).Decode(&result); err == nil {
		res = append(res, FieldValidationErrorResponse{FieldName: "name", Message: "This permission is already registered."})
	}

	return res
}

func (p *Permission) CreatePermission() error {
	_, err := permissionCollection.InsertOne(context.Background(), p)

	if err != nil {
		return err
	}

	return nil
}

func AllPermissions(isActive *bool) ([]PermissionList, error) {
	var permissions []PermissionList

	ctx := context.Background()
	filter := bson.M{}
	if isActive != nil {
		filter["isactive"] = *isActive
	}

	cursor, err := permissionCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var permission PermissionList
		if err := cursor.Decode(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func GetPermissionByID(id primitive.ObjectID, checkActiveness bool) (Permission, bool) {
	var permission Permission
	query := bson.M{"_id": id}

	err := permissionCollection.FindOne(context.Background(), query).Decode(&permission)
	if err != nil {
		return Permission{}, false
	}

	if checkActiveness && !permission.IsActive {
		return Permission{}, false
	}

	return permission, true
}
