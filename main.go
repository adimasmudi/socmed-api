package main

import (
	"socmed-api/configs"
	"socmed-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// run database
	configs.ConnectDB()

	// add routes
	routes.PostRoute(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from fiber and mongo"})
	})

	app.Listen(":6000")
}
