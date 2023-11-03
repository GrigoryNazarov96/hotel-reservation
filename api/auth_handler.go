package api

import (
	"errors"
	"fmt"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store *db.Store
}

func NewAuthHandler(s *db.Store) *AuthHandler {
	return &AuthHandler{
		store: s,
	}
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var dto types.LoginDTO
	if err := c.BodyParser(&dto); err != nil {
		return nil
	}
	user, err := h.store.User.GetUserByEmail(c.Context(), dto.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncPwd), []byte(dto.Pwd)); err != nil {
		return fmt.Errorf("invalid credentials")
	}
	return c.JSON(map[string]string{"message": "logged in successfully"})
}
