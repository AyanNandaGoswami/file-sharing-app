package models

import (
	"context"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/nexus/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var serviceCollection *mongo.Collection

const serviceCollectionName = "services"

func init() {
	serviceCollection = database.DB.Collection(serviceCollectionName)
}

type Service struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UUID        string             `json:"uuid"`
	Name        string             `json:"name" validate:"required=This field is required."`
	Code        string             `json:"code" validate:"omitempty"`
	Description string             `json:"description" validate:"omitempty"`
}

func (s *Service) ValidateServiceRegistrationPayload() []FielValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(s)
	var res []FielValidationErrorResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, FielValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	// validate email already exists or not
	var result Service
	query := bson.M{"name": s.Name}

	if err := serviceCollection.FindOne(context.Background(), query).Decode(&result); err == nil {
		res = append(res, FielValidationErrorResponse{FieldName: "name", Message: "This service is already registered."})
	}

	return res
}

func (s *Service) RegisterService() error {

	s.UUID = uuid.New().String()
	_, err := serviceCollection.InsertOne(context.Background(), s)

	if err != nil {
		return err
	}
	return nil
}

func GetAllServices() ([]Service, error) {
	// Define a slice to store services
	var services []Service

	// Define a filter to match all documents
	filter := bson.D{}

	// Find all services
	cursor, err := serviceCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and decode each document
	for cursor.Next(context.Background()) {
		var service Service
		if err := cursor.Decode(&service); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func GetServiceById(serviceID string) (Service, error) {
	var service Service

	objectID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return Service{}, err
	}

	filterQuery := bson.M{"_id": objectID}

	if err := serviceCollection.FindOne(context.Background(), filterQuery).Decode(&service); err != nil {
		return Service{}, err
	}

	return service, nil

}
