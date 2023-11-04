package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"github.com/gofiber/fiber/v2"
)

func makeTestUser(t *testing.T, s *db.Store) *types.User {
	user, err := types.NewUserFromDTO(types.CreateUserDTO{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john_doe@test.com",
		Pwd:       "test1234",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.User.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestLoginSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	tu := makeTestUser(t, tdb.store)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store)
	app.Post("/", authHandler.HandleLogin)

	loginDto := types.LoginDTO{
		Email: "john_doe@test.com",
		Pwd:   "test1234",
	}

	b, _ := json.Marshal(loginDto)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	res, _ := app.Test(req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status to be 200, got %d", res.StatusCode)
	}

	var r types.LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Error(err)
	}
	if r.User.Email != loginDto.Email {
		t.Errorf("expected email to be %s, but got %s", loginDto.Email, r.User.Email)
	}
	if len(r.Token) == 0 {
		t.Errorf("expected JWT to be presented")
	}
	tu.EncPwd = ""
	if !reflect.DeepEqual(tu, r.User) {
		t.Errorf("expected the user remain the same")
	}
}

func TestWrongPwdLogin(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	makeTestUser(t, tdb.store)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store)
	app.Post("/", authHandler.HandleLogin)

	loginDto := types.LoginDTO{
		Email: "john_doe@test.com",
		Pwd:   "tezt1234",
	}

	b, _ := json.Marshal(loginDto)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	res, _ := app.Test(req)
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status to be 400, got %d", res.StatusCode)
	}

	var r map[string]string
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Fatal(err)
	}
	if r["message"] != "invalid credentials" {
		t.Errorf("expected to get 'invalid credentials' message, got %s", r["message"])
	}
}
