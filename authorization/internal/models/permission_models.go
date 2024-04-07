package models

import (
	"context"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var permissionCollection *mongo.Collection

const permissionCollectionName = "permissions"

func init() {
	permissionCollection = database.DB.Collection(permissionCollectionName)
}

type Permission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	APIEndpoint string             `json:"api_endpoint"`
	Service     primitive.ObjectID `bson:"service"`
	IsActve     bool               `json:"is_active"`
}

func (p *Permission) CreatePermission() error {
	_, err := permissionCollection.InsertOne(context.Background(), p)

	if err != nil {
		return err
	}

	return nil
}
