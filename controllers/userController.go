package controllers

import (
	"context"
	"net/http"
	"socmed-api/configs"
	"socmed-api/models"
	"socmed-api/responses"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User

	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// validate the request body
	if validationError := validate.Struct(&user); validationError != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationError.Error()}})
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username, "email": user.Email})
	if count >= 1 {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "email or username already exist", Data: &fiber.Map{"data": count}})
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Username: user.Username,
		Email:    user.Email,
		Password: string(hashedPassword),
	}

	result, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})

}

// func GetAUser() {
// 	return
// }

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	err := userCollection.FindOne(ctx, bson.M{"email": payload.Email}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	error := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "Wrong Password", Data: &fiber.Map{"data": error.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PostResponse{Status: http.StatusOK, Message: "success, ada", Data: &fiber.Map{"data": user}})

}

// func Logout() {
// 	return
// }
