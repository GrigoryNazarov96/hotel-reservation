package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(s *db.Store) *UserHandler {
	return &UserHandler{
		store: s,
	}
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var dto types.CreateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}
	if errors := dto.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromDTO(dto)
	if err != nil {
		return err
	}
	newUser, err := h.store.User.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(newUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	u, err := h.store.User.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NewError(http.StatusNotFound, "resource not found")
		}
		return err
	}
	return c.JSON(u)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u, err := h.store.User.GetUsers(c.Context())
	if err != nil {
		return NewError(http.StatusNotFound, "resource not found")
	}
	return c.JSON(u)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.store.User.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"message": fmt.Sprintf("User %s deleted", userID)})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		dto    types.UpdateUserDTO
		userID = c.Params("id")
	)
	if err := c.BodyParser(&dto); err != nil {
		return err
	}
	_, err := h.store.User.UpdateUser(c.Context(), userID, dto)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"message": fmt.Sprintf("user %s updated successfully", userID)})
}
