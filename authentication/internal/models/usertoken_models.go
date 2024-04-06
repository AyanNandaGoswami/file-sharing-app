package models

import (
	"context"
	"time"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var tokenCollection *mongo.Collection

const tokenCollectionName = "blacklistedtokens"

func init() {
	tokenCollection = database.DB.Collection(tokenCollectionName)
}

// Token is the model for blacklisted
type BlacklistedToken struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBlacklistedToken(tkn string) error {
	token := BlacklistedToken{Token: tkn, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	_, err := tokenCollection.InsertOne(context.Background(), token)
	if err != nil {
		return err
	}
	return nil
}

func IsTokenBlacklisted(tkn string) bool {
	var token BlacklistedToken
	query := bson.M{"token": tkn}

	err := tokenCollection.FindOne(context.Background(), query).Decode(&token)

	return err == nil

}
