package api

import (
	"net/http"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(s *db.Store) *HotelHandler {
	return &HotelHandler{
		store: s,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	filter := bson.M{}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter)
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NewError(http.StatusBadRequest, "invalid id provided")
	}
	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	return c.JSON(rooms)
}
