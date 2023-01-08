package routes

import (
	"socmed-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func PostRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Post("/post", controllers.CreatePost)
	app.Get("/post/:postId", controllers.GetAPost)
	app.Put("/post/:postId", controllers.EditAPost)
	app.Delete("/post/:postId", controllers.DeleteAPost)
	app.Get("/posts", controllers.GetAllPosts)
}
