package main

import (
	"context"
	"fmt"
	"log"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        = context.Background()
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
)

func seedHotels(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Type:      types.Single,
			BasePrice: 184.9,
		},
		{
			Type:      types.Double,
			BasePrice: 349.9,
		},
		{
			Type:      types.Suite,
			BasePrice: 489.9,
		},
		{
			Type:      types.Deluxe,
			BasePrice: 1449.9,
		},
	}
	createdHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = createdHotel.Id
		createdRoom, err := roomStore.CreateRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(createdRoom)
	}
	fmt.Println(createdHotel)
}

func main() {
	seedHotels("White Lotus", "Palermo", 5)
}

func init() {
	_, DB_NAME, DB_URI := db.Init()
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, DB_NAME)
	roomStore = db.NewMongoRoomStore(client, hotelStore, DB_NAME)
}
