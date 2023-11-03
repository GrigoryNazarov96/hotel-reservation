package main

import (
	"context"
	"flag"
	"log"

	"github.com/GrigoryNazarov96/hotel-reservation.git/api"
	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
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
	userStore := db.NewMongoUserStore(client)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)

	//handlers init
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(hotelStore, roomStore)

	//app init
	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	//users
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Patch("/user/:id", userHandler.HandleUpdateUser)

	//hotels
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotel)

	//listener
	app.Listen(*listenAddr)
}
