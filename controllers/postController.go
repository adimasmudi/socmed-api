package controllers

import (
	"context"
	"net/http"
	"socmed-api/configs"
	"socmed-api/models"
	"socmed-api/responses"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts")
var validate = validator.New()

func CreatePost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var post models.Post

	defer cancel()

	//validate the request body
	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// validate the request body
	if validationError := validate.Struct(&post); validationError != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationError.Error()}})
	}

	newPost := models.Post{
		Id:      primitive.NewObjectID(),
		Owner:   post.Owner,
		Title:   post.Title,
		Content: post.Content,
	}

	result, err := postCollection.InsertOne(ctx, newPost)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.PostResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	postId := c.Params("postId")
	var post models.Post
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	err := postCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&post)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PostResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": post}})
}

func EditAPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	postId := c.Params("postId")
	var post models.Post
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	// validate request body
	if err := c.BodyParser(&post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// use validator library to validate required field
	if validationErr := validate.Struct(&post); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PostResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"owner": post.Owner, "title": post.Title, "content": post.Content}

	result, err := postCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated post details
	var updatedPost models.Post
	if result.MatchedCount == 1 {
		err := postCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPost)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.PostResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedPost}})

}

func DeleteAPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	postId := c.Params("postId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	result, err := postCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.PostResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.PostResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllPosts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.Post
	defer cancel()

	results, err := postCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.Post
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PostResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.PostResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}
