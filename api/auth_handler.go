package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
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
		return err
	}
	user, err := h.store.User.GetUserByEmail(c.Context(), dto.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return NewError(http.StatusBadRequest, "invalid credentials")
		}
		return err
	}
	if !types.IsValidPassword(user.EncPwd, dto.Pwd) {
		return NewError(http.StatusBadRequest, "invalid credentials")
	}
	res := types.LoginResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}
	return c.JSON(res)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.Id,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token ", err)
	}
	return tokenStr
}
