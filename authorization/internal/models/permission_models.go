package models

import (
	"context"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetAPIEndpointsFromPermissions(permissionPrimitiveIds []primitive.ObjectID) (map[string]string, error) {
	// Step 1: Retrieve the Permission documents using the Permission IDs
	var permissions []Permission
	permissionQuery := bson.M{"_id": bson.M{"$in": permissionPrimitiveIds}}

	// Assume getCollection is a function that returns the correct collection
	collection := getCollection(permissionCollection)
	cursor, err := collection.Find(context.Background(), permissionQuery)
	if err != nil {
		return nil, err // Return the error if permissions cannot be fetched
	}
	defer cursor.Close(context.Background()) // Always close the cursor once done

	// Step 2: Decode the results into the permissions slice
	if err := cursor.All(context.Background(), &permissions); err != nil {
		return nil, err // Return the error if decoding fails
	}

	// Step 3: Retrieve the associated API Endpoints for each permission
	var allApis []APIEndpoints // Collect all API endpoints in a slice

	// Getting the API Endpoints collection
	apiEndpointCollection := getCollection(apiEndpointCollection)

	for _, permission := range permissions {
		// Fetch API endpoints for each permission, using the ObjectIDs from permission.APIs
		cursor, err := apiEndpointCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": permission.APIs}})
		if err != nil {
			return nil, err // Return the error if APIs cannot be fetched
		}
		defer cursor.Close(context.Background()) // Always close the cursor once done

		// Decode the API results into the apis slice
		var apis []APIEndpoints
		if err := cursor.All(context.Background(), &apis); err != nil {
			return nil, err // Return the error if decoding fails
		}

		// Add the fetched APIs to the allApis slice
		allApis = append(allApis, apis...)
	}

	// Step 4: Map the API Endpoints to a map (URL -> Method)
	apiEndpoints := make(map[string]string) // Create a map to store URL and Method

	// Now loop through the allApis slice outside the permission loop
	for _, api := range allApis {
		// Add each API URL as the key, and its method as the value in the map
		apiEndpoints[api.URL] = api.Method
	}

	// Return the collected API URLs and Methods
	return apiEndpoints, nil
}
