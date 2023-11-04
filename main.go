package main

import (
	"context"
	"flag"
	"log"

	"github.com/GrigoryNazarov96/hotel-reservation.git/api"
	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/middleware"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"Error": err.Error()})
	},
}

func main() {
	//env variable analog (to change the port: make run --listenAddr :8080)
	listenAddr := flag.String("listenAddr", ":5020", "The listen address for API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	//stores init
	userStore := db.NewMongoUserStore(client, db.DB_NAME)
	hotelStore := db.NewMongoHotelStore(client, db.DB_NAME)
	roomStore := db.NewMongoRoomStore(client, hotelStore, db.DB_NAME)
	store := &db.Store{
		User:  userStore,
		Hotel: hotelStore,
		Room:  roomStore,
	}

	//handlers init
	userHandler := api.NewUserHandler(store)
	hotelHandler := api.NewHotelHandler(store)
	authHandler := api.NewAuthHandler(store)
	bookingHandler := api.NewBookingHandler(store)
	//app init
	app := fiber.New(config)

	//API grouping
	apiv1 := app.Group("/api/v1", middleware.JWTAuth(userStore))
	auth := app.Group("/api/")

	//users
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Patch("/user/:id", userHandler.HandleUpdateUser)

	//auth
	auth.Post("/login", authHandler.HandleLogin)

	//hotels
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	//bookings
	apiv1.Post("/room/:id/book", bookingHandler.HandleBookRoom)

	//listener
	app.Listen(*listenAddr)
}
