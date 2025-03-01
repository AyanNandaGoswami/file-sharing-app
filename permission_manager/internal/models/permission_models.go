package models

import (
	"context"
	"fmt"

	common_models "github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/microservices/file-sharing-app/permission_manager/internal/database"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// declare colleciton constants
const permissionCollection = "permissions"
const apiEndpointCollection = "api_endpoints"

var dbCollections = make(map[string]*mongo.Collection)
var collectionNames = []string{permissionCollection, apiEndpointCollection}

func init() {

	for _, name := range collectionNames {
		dbCollections[name] = database.DB.Collection(name)
	}
}

func getCollection(collectionName string) *mongo.Collection {
	return dbCollections[collectionName]
}

type APIEndpoints struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `json:"name" validate:"required=This field is required."`
	URL    string             `json:"url" validate:"required=This field is required."`
	Method string             `json:"method" validate:"required=This field is required."`
}

type Permission struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"`
	Name     string               `json:"name" validate:"required=This field is required."`
	IsActive bool                 `json:"is_active" validate:"required=This field is required."`
	APIs     []primitive.ObjectID `json:"apis" bson:"apis" validate:"required=This field is required."`
}

type PermissionList struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" validate:"required=This field is required."`
	IsActive bool               `json:"is_active" validate:"required=This field is required."`
}

// APIEndpoints
func (apis *APIEndpoints) ValidateAPIEndpointsRegistrationPayload() []common_models.FieldValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(apis)
	var res []common_models.FieldValidationErrorResponse

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, common_models.FieldValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	// validate endpoint is already registered or not with the same method
	var result APIEndpoints
	query := bson.M{"url": apis.URL, "method": apis.Method}

	// Access the collection dynamically
	collection := getCollection(apiEndpointCollection)

	if err := collection.FindOne(context.Background(), query).Decode(&result); err == nil {
		res = append(res, common_models.FieldValidationErrorResponse{
			FieldName: "url", Message: fmt.Sprintf("API endpoint %s is already registered with the same method %s", apis.URL, apis.Method)})
	}

	return res
}

func (apis *APIEndpoints) RegisterNewAPIEndpoint() error {
	collection := getCollection(apiEndpointCollection)
	_, err := collection.InsertOne(context.Background(), apis)
	return err
}

func AllAPIEndpoints() ([]APIEndpoints, error) {
	var endpoints []APIEndpoints
	// Explicitly define context
	ctx := context.Background()

	// Access the collection dynamically
	collection := getCollection(apiEndpointCollection)

	// Define an empty filter (bson.M{}) to fetch all documents
	cursor, err := collection.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document
	for cursor.Next(ctx) {
		var endpoint APIEndpoints
		if err := cursor.Decode(&endpoint); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}

	// Check if there was an error with the cursor
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return endpoints, nil
}

// Permissions
func (p *Permission) ValidatePermissionRegistrationPayload() []common_models.FieldValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(p)
	var res []common_models.FieldValidationErrorResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, common_models.FieldValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	// validate permission is already exists or not
	var result Permission
	query := bson.M{"name": p.Name}

	// Access the collection dynamically
	collection := getCollection(permissionCollection)

	if err := collection.FindOne(context.Background(), query).Decode(&result); err == nil {
		res = append(res, common_models.FieldValidationErrorResponse{FieldName: "name", Message: "This permission is already registered."})
	}

	return res
}

func (p *Permission) CreatePermission() error {
	// Access the collection dynamically
	collection := getCollection(permissionCollection)

	_, err := collection.InsertOne(context.Background(), p)

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

	// Access the collection dynamically
	collection := getCollection(permissionCollection)
	cursor, err := collection.Find(ctx, filter)
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

func GetPermissionByID(id primitive.ObjectID, checkActiveness bool) *Permission {
	var permission Permission
	query := bson.M{"_id": id}

	// Access the collection dynamically
	collection := getCollection(permissionCollection)

	// Find the permission by ID
	err := collection.FindOne(context.Background(), query).Decode(&permission)
	if err != nil {
		// If there's an error (e.g., permission not found), return nil
		return nil
	}

	// If checkActiveness is true and the permission is not active, return nil
	if checkActiveness && !permission.IsActive {
		return nil
	}

	// Return the permission if it's found and active
	return &permission
}
