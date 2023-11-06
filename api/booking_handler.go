package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
		return NewError(http.StatusBadRequest, "invalid id provided")
	}
	var dto types.BookingDTO
	if err := c.BodyParser(&dto); err != nil {
		return err
	}
	if err := dto.Validate(); err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}
	user, err := getAuthenticatedUser(c.Context())
	if err != nil {
		return err
	}
	ok, err := h.isAvailable(c.Context(), roomId, dto)
	if err != nil {
		return err
	}
	if !ok {
		return c.JSON(map[string]string{"message": "room is not available for those dates"})
	}
	booking := &types.Booking{
		UserId:   user.Id,
		RoomId:   roomId,
		FromDate: dto.FromDate,
		TillDate: dto.TillDate,
		Guests:   dto.Guests,
	}
	newBooking, err := h.store.Booking.CreateBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(newBooking)
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	user, err := getAuthenticatedUser(c.Context())
	if err != nil {
		return err
	}
	filter := bson.M{"userId": user.Id}
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := getAuthenticatedUser(c.Context())
	if err != nil {
		return err
	}
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	if user.Id != booking.UserId {
		return NewError(http.StatusUnauthorized, "you can browse only your bookings")
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleUpdateBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	user, err := getAuthenticatedUser(c.Context())
	if err != nil {
		return err
	}
	if user.Id != booking.UserId {
		return NewError(http.StatusUnauthorized, "you can update only your bookings")
	}
	var dto types.UpdateBookingDTO
	if err := c.BodyParser(&dto); err != nil {
		return err
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), id, dto); err != nil {
		return err
	}
	return c.JSON(map[string]string{"message": fmt.Sprintf("Booking %s updated successfully", id)})
}

func (h *BookingHandler) isAvailable(c context.Context, roomId primitive.ObjectID, dto types.BookingDTO) (bool, error) {
	filter := bson.M{
		"fromDate": bson.M{
			"$gte": dto.FromDate,
		},
		"tillDate": bson.M{
			"$lte": dto.FromDate,
		},
		"roomId": roomId,
	}
	bookings, err := h.store.Booking.GetBookings(c, filter)
	if err != nil {
		return false, NewError(http.StatusBadRequest, err.Error())
	}
	return len(bookings) == 0, nil
}
