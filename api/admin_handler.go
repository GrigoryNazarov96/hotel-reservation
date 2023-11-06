package api

import (
	"net/http"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type AdminHandler struct {
	store *db.Store
}

func NewAdminHandler(s *db.Store) *AdminHandler {
	return &AdminHandler{
		store: s,
	}
}

func (h *AdminHandler) HandleGetBookings(c *fiber.Ctx) error {
	filter := bson.M{}
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	return c.JSON(bookings)
}
