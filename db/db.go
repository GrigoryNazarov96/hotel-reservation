package db

import "context"

const (
	DB_NAME      string = "hotel-reservation"
	DB_URI       string = "mongodb://localhost:27017"
	TEST_DB_NAME string = "test_db"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
type Dropper interface {
	Drop(context.Context) error
}
