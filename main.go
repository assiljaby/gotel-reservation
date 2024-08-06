package main

import "github.com/gofiber/fiber/v3"

func main() {
	app := fiber.New()
	app.Get("/foo", handleFoo)
	app.Listen(":3000")
}

func handleFoo(c fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Bar!"})
}