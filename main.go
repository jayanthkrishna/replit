package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/repl/:id", controllers.createrepl)

	app.Listen(":3000")
}
