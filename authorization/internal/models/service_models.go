package models

import (
	"context"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/database"
	"github.com/go-playground/validator/v10"
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
	_, err := serviceCollection.InsertOne(context.Background(), s)

	if err != nil {
		return err
	}
	return nil
}
