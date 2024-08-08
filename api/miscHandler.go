package api

import "github.com/gofiber/fiber/v2"

func HandleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Bar!"})
}