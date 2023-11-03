package api

import (
	"errors"
	"fmt"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var dto types.CreateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return err
	}
	if errors := dto.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromDTO(dto)
	if err != nil {
		return err
	}
	newUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(newUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	u, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "Not Found"})
		}
		return err
	}
	return c.JSON(u)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(u)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
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
	_, err := h.userStore.UpdateUser(c.Context(), userID, dto)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"message": fmt.Sprintf("user %s updated successfully", userID)})
}
