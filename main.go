package main

import (
	"flag"

	"github.com/assiljaby/gotel-reservation/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenPort := flag.String("listenPort", ":3000", "The server is listening to this port")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	app.Get("/foo", handleFoo)
	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/users/:id", api.HandleGetUser)

	app.Listen(*listenPort)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "Bar!"})
}