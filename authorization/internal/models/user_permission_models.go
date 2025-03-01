package models

import (
	"context"
	"fmt"
	"log"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/database"

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
