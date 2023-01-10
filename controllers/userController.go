package controllers

import (
	"context"
	"fmt"
	"net/http"
	"socmed-api/configs"
	"socmed-api/models"
	"socmed-api/responses"
	"time"

	"github.com/golang-jwt/jwt/v4"

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

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
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

	fmt.Println(user.Username, user.Password)

	err := userCollection.FindOne(ctx, bson.M{"username": user.Username, "password": user.Password}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PostResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user, "token": t}})
}

// func Logout() {
// 	return
// }
