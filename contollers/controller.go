package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	Id          string `json:"id"`
	Environment string `json:"environment"`
}

type Response struct {
	Url string `json:"url"`
}

func (k *Kconfig) Createrepl(c *fiber.Ctx) error {

	var request Request

	err := c.BodyParser(&request)

	if err != nil {
		return err
	}

	url, err := k.createResources(request.Id, request.Environment, "default")

	if err != nil {
		fmt.Printf("error while creating resources: %s\n", err.Error())
		return err
	}

	return c.JSON(Response{
		Url: url,
	})

}
