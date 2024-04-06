package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MyClient *mongo.Client
var DB *mongo.Database
var databaseName = "auth_service_database"
var databaseCollections = []string{"users", "blacklistedtokens"}

func createDdatabaseCollections(ctx context.Context) {
	for _, collectionName := range databaseCollections {
		err := DB.CreateCollection(ctx, collectionName)
		if err != nil {
			log.Fatalf("[error] getting error during creation for collection => %s. [errorDetail] %s", collectionName, err.Error())
		}
	}
}

func createOrRetrieveDatabase(ctx context.Context) {
	DB = MyClient.Database(databaseName)
	createDdatabaseCollections(ctx)
}

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://root:root@localhost:27020")

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
