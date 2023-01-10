package main

import (
	"socmed-api/configs"
	"socmed-api/routes"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	app := fiber.New()

	// run database
	configs.ConnectDB()

	// add routes
	routes.PostRoute(app)
	routes.UserRoute(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from fiber and mongo"})
	})

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	app.Listen(":6000")
}
