package api

import (
	"github.com/assiljaby/gotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("ScrubLord")
}

func HandleGetUsers(c *fiber.Ctx) error {
	scrubLoard := types.User{
		FirstName: "Scrub",
		LastName: "Lord",
	}

	return c.JSON(scrubLoard)
}
