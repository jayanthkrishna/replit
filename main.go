package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	controllers "github.com/jayanthkrishna/replit/contollers"
)

func main() {
	clientset := controllers.KubernetesConfig()
	kconfig := controllers.NewKconfig(clientset)

	app := fiber.New()

	app.Post("/repl", kconfig.Createrepl)

	fmt.Println("Server running on :3000")
	app.Listen(":3000")
}
