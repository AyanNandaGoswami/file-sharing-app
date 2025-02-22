package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MyClient *mongo.Client
var DB *mongo.Database
var databaseName = "authorization_service_database"
var databaseCollections = []string{"permissions", "userpermissions", "api_endpoints"}

func createDatabaseCollections(ctx context.Context) {
	// List all collections in the database once
	cursor, err := DB.ListCollections(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}
	defer cursor.Close(ctx)

	// Create a map to store existing collection names
	existingCollections := make(map[string]bool)

	// Populate the map with the collection names
	for cursor.Next(ctx) {
		var collection struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&collection); err != nil {
			log.Fatalf("Error decoding collection name: %v", err)
		}
		existingCollections[collection.Name] = true
	}

	// Now check if each collection from the list should be created
	for _, collectionName := range databaseCollections {
		if _, exists := existingCollections[collectionName]; !exists {
			// If collection doesn't exist, create it
			if err := DB.CreateCollection(ctx, collectionName); err != nil {
				log.Fatalf("Error creating collection %s: %v", collectionName, err)
			}
		}
	}
}

func createOrRetrieveDatabase(ctx context.Context) {
	DB = MyClient.Database(databaseName)
	createDatabaseCollections(ctx)
}

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27020")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Assign the client and database to package-level variables
	MyClient = client

	createOrRetrieveDatabase(ctx)
	fmt.Println("Database is connected successfully and ready to use.")
}
