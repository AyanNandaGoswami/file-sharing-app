package models

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

const userCollectionName = "users"

type CreateUser struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Firstname  string             `json:"firstname" validate:"required=This field is required."`
	Lastname   string             `json:"lastname" validate:"required=This field is required."`
	Middlename string             `json:"middlename" validate:"omitempty"`
	Email      string             `json:"email" validate:"required=This field is required.,email=This email address is not valid."`
	Password   string             `json:"password" validate:"required=This field is required."`
	UUID       string             `json:"uuid"`
	IsActive   bool               `json:"is_active"`
	IsVerified bool               `json:"is_verified"`
	CreatedOn  time.Time          `json:"created_on"`
	UpdatedOn  time.Time          `json:"updated_on"`
	LastLogin  time.Time          `json:"last_login"`
	IsAdmin    bool               `json:"is_admin" validate:"omitempty"`
}

type UserDetail struct {
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email"`
	UUID       string `json:"uuid"`
}

type UserUpdate struct {
	Firstname  string `json:"firstname,omitempty"`
	Lastname   string `json:"lastname,omitempty"`
	Middlename string `json:"middlename,omitempty"`
	Email      string `json:"email,omitempty"`
}

type Login struct {
	Email    string `json:"email" validate:"required=This field is required.,email=This email address is not valid."`
	Password string `json:"password" validate:"required=This field is required."`
}

func init() {
	userCollection = database.DB.Collection(userCollectionName)
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

	if err := userCollection.FindOne(context.Background(), query).Decode(&result); err == nil {
		res = append(res, FielValidationErrorResponse{FieldName: "email", Message: "This email already taken."})
	}

	return res
}

func (u *UserUpdate) ValidateUserUpdatePayload() []FielValidationErrorResponse {
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

	if err := userCollection.FindOne(context.Background(), query).Decode(&result); err == nil {
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

	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

// update user
func (u *UserUpdate) UpdateUserByID(uuid string) error {
	update := bson.M{}

	v := reflect.ValueOf(*u)

	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldValue := v.Field(i).Interface().(string)

		if fieldValue != "" {
			update[fieldName] = fieldValue
		}
	}

	if len(update) == 0 {
		return errors.New("no fields to be updated")
	}

	filter := bson.M{"uuid": uuid}
	_, err := userCollection.UpdateOne(context.Background(), filter, bson.M{"$set": update})

	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(uuid string) error {
	query := bson.M{"uuid": uuid}

	_, err := userCollection.DeleteOne(context.Background(), query)

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

	if err := userCollection.FindOne(context.Background(), query).Decode(&user); err != nil {
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

	if err := userCollection.FindOne(context.Background(), query).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}
