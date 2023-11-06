package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Init() (string, string, string) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		TEST_DB_NAME string = os.Getenv("TEST_DB_NAME")
		DB_NAME      string = os.Getenv("DB_NAME")
		DB_URI       string = os.Getenv("DB_URI")
	)
	return TEST_DB_NAME, DB_NAME, DB_URI
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
type Dropper interface {
	Drop(context.Context) error
}
