package models

import (
	"context"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authorization/internal/database"
	"github.com/go-playground/validator/v10"
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
	Code        string             `json:"code" validate:"required=This field is required."`
	Description string             `json:"description" validate:"omitempty"`
}

type Permission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	APIEndpoint string             `json:"api_endpoint"`
	service     primitive.ObjectID `bson:"service_id"`
	IsActve     bool               `json:"is_active"`
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

	// // validate email already exists or not
	// var result CreateUser
	// query := bson.M{"email": u.Email}

	// if err := userCollection.FindOne(context.Background(), query).Decode(&result); err == nil {
	// 	res = append(res, FielValidationErrorResponse{FieldName: "email", Message: "This email already taken."})
	// }

	return res
}

func (s *Service) RegisterService() error {
	_, err := serviceCollection.InsertOne(context.Background(), s)

	if err != nil {
		return err
	}
	return nil
}
