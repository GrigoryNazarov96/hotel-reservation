package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sample_db struct {
	store *db.Store
}

func setup(t *testing.T) *sample_db {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}
	return &sample_db{
		store: &db.Store{
			User: db.NewMongoUserStore(client, db.TEST_DB_NAME),
		},
	}
}

func (s *sample_db) teardown(t *testing.T) error {
	return s.store.User.Drop(context.TODO())
}

func TestCreateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	userHandler := NewUserHandler(tdb.store)

	app := fiber.New()
	app.Post("/", userHandler.HandleCreateUser)

	dto := types.CreateUserDTO{
		Email:     "test@test.com",
		FirstName: "John",
		LastName:  "Doe",
		Pwd:       "test1234",
	}

	b, _ := json.Marshal(dto)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	res, _ := app.Test(req)

	var user types.User
	json.NewDecoder(res.Body).Decode(&user)

	if user.FirstName != dto.FirstName {
		t.Errorf("expected firstname to be %s, but got %s", dto.FirstName, user.FirstName)
	}
	if user.LastName != dto.LastName {
		t.Errorf("expected lastname to be %s, but got %s", dto.LastName, user.LastName)
	}
	if user.Email != dto.Email {
		t.Errorf("expected email to be %s, but got %s", dto.Email, user.Email)
	}
}
