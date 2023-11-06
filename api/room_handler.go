package api

import (
	"net/http"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(s *db.Store) *RoomHandler {
	return &RoomHandler{
		store: s,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	return c.JSON(rooms)
}
