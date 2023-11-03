package api

import (
	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	roomStore  db.RoomStore
	hotelStore db.HotelStore
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qp HotelQueryParams
	if err := c.QueryParser(&qp); err != nil {
		return err
	}
	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}
