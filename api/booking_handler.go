package api

import (
	"net/http"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(s *db.Store) *BookingHandler {
	return &BookingHandler{
		store: s,
	}
}

func (h *BookingHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	var dto types.BookingDTO
	if err := c.BodyParser(&dto); err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.SendStatus(http.StatusInternalServerError)
	}
	booking := types.Booking{
		UserId:   user.Id,
		RoomId:   roomId,
		FromDate: dto.FromDate,
		TillDate: dto.TillDate,
		Guests:   dto.Guests,
	}
	return c.JSON(booking)
}
