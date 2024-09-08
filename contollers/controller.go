package controllers

import "github.com/gofiber/fiber/v2"

type Request struct {
	Id          string `json: "id"`
	environment string `json:"environment"`
}

func createrepl(c *fiber.Ctx) error {

	var request Request

	err := c.BodyParser(&request)

	if err != nil {
		return err
	}

	return nil

}
