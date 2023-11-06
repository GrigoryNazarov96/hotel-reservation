package main

import (
	"context"
	"log"
	"os"

	"github.com/GrigoryNazarov96/hotel-reservation.git/api"
	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	_, DB_NAME, DB_URI := db.Init()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	//stores init
	userStore := db.NewMongoUserStore(client, DB_NAME)
	hotelStore := db.NewMongoHotelStore(client, DB_NAME)
	roomStore := db.NewMongoRoomStore(client, hotelStore, DB_NAME)
	bookingStore := db.NewMongoBookingStore(client, DB_NAME)
	store := &db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	//handlers init
	userHandler := api.NewUserHandler(store)
	hotelHandler := api.NewHotelHandler(store)
	authHandler := api.NewAuthHandler(store)
	bookingHandler := api.NewBookingHandler(store)
	roomHandler := api.NewRoomHandler(store)
	adminHandler := api.NewAdminHandler(store)

	//app init
	app := fiber.New(config)

	//API grouping
	apiv1 := app.Group("/api/v1", middleware.JWTAuth(userStore))
	auth := app.Group("/api")
	admin := apiv1.Group("admin", middleware.AdminAuth)

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

	//rooms
	apiv1.Get("/room", roomHandler.HandleGetRooms)

	//bookings
	apiv1.Post("/room/:id/book", bookingHandler.HandleBookRoom)
	apiv1.Get("/booking", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Patch("/booking/:id", bookingHandler.HandleUpdateBooking)

	//admin
	admin.Get("/booking", adminHandler.HandleGetBookings)

	//listener
	port := os.Getenv("PORT")
	app.Listen(port)
}
