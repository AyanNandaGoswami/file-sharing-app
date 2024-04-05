package models

import (
	"context"
	"time"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection

const collectionName = "users"

type CreateUser struct {
	Firstname  string    `json:"firstname" validate:"required=This field is required."`
	Lastname   string    `json:"lastname" validate:"required=This field is required."`
	Middlename string    `json:"middlename" validate:"omitempty"`
	Email      string    `json:"email" validate:"required=This field is required.,email=This email address is not valid."`
	Password   string    `json:"password" validate:"required=This field is required."`
	UUID       string    `json:"uuid"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	CreatedOn  time.Time `json:"created_on"`
	UpdatedOn  time.Time `json:"updated_on"`
	LastLogin  time.Time `json:"last_login"`
}

type UserDetail struct {
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email"`
	UUID       string `json:"uuid"`
}

type Login struct {
	Email    string `json:"email" validate:"required=This field is required.,email=This email address is not valid."`
	Password string `json:"password" validate:"required=This field is required."`
}

func init() {
	collection = database.DB.Collection(collectionName)
}

// user registration
func (u *CreateUser) ValidateUserRegistrationPayload() []FielValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(u)
	var res []FielValidationErrorResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, FielValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	// validate email already exists or not
	var result CreateUser
	query := bson.M{"email": u.Email}

	if err := collection.FindOne(context.Background(), query).Decode(&result); err == nil {
		res = append(res, FielValidationErrorResponse{FieldName: "email", Message: "This email already taken."})
	}

	return res
}

func (user *CreateUser) CreateNewUser() error {

	// make password hashed
	hashed_password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	user.Password = string(hashed_password)
	user.CreatedOn = time.Now()
	user.UpdatedOn = time.Now()
	user.IsActive = true
	user.IsVerified = true
	user.UUID = uuid.New().String()

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

// login
func (l *Login) ValiadteLoginPayload() []FielValidationErrorResponse {
	validate := validator.New()
	err := validate.Struct(l)
	var res []FielValidationErrorResponse

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			res = append(res, FielValidationErrorResponse{FieldName: err.StructField(), Message: err.Param()})
		}
	}

	return res
}

func (l *Login) Authenticate() (string, error) {
	// all logic for Login
	var user CreateUser
	query := bson.M{"email": l.Email}

	if err := collection.FindOne(context.Background(), query).Decode(&user); err != nil {
		return "", err
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.Password))
		if err != nil {
			return "", err
		}
	}
	return user.UUID, nil
}

func GetUserByUUID(userId string) (UserDetail, error) {
	var result UserDetail
	query := bson.M{"uuid": userId}

	if err := collection.FindOne(context.Background(), query).Decode(&result); err != nil {
		return result, err
	}
	return result, nil

}
