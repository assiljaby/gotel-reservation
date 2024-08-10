package api

import (
	"errors"

	"github.com/assiljaby/gotel-reservation/db"
	"github.com/assiljaby/gotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.UserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	errors := params.Validate()
	if len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	res, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(res)
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (	
		id = c.Params("id")
		userPrms types.UserParams
	)
	if err := c.BodyParser(&userPrms); err != nil {
		return err
	}

	errors := userPrms.Validate()
	if len(errors) > 0 {
		return c.JSON(errors)
	}


	updatedUser, err := types.NewUserFromParams(userPrms)
	if err != nil {
		return err
	}


	if err = h.userStore.UpdateUser(c.Context(), id, *updatedUser); err != nil {
		return err 
	}
	
	return c.JSON(map[string]string{"msg": "user updated"})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(map[string]string{"error": "user does not exist"})
		}
		return err
	}
	return c.JSON(map[string]string{"msg": "User Deleted Successfully"})
}