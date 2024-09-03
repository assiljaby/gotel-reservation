package api

import (
	"net/http"
	"time"

	"github.com/assiljaby/gotel-reservation/db"
	"github.com/assiljaby/gotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate  time.Time `json:"fromDate"`
	TillDate  time.Time `json:"tillDate"`
	NumPerson int       `json:"numPerson"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
	}

	_ = booking
	return nil
}
