package routes

import (
	"socmed-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Post("/user/register", controllers.Register)
	// app.Get("/user/:userId", controllers.GetAUser)
	app.Post("/user/login", controllers.Login)
	// app.Post("/logout", controllers.Logout)
}
