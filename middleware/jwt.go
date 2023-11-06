package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/GrigoryNazarov96/hotel-reservation.git/api"
	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(s db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Bearer"]
		if !ok {
			return api.NewError(http.StatusUnauthorized, "unauthorized")
		}
		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return api.NewError(http.StatusUnauthorized, "please login again")
		}
		userId := claims["id"].(string)
		user, err := s.GetUserByID(c.Context(), userId)
		if err != nil {
			return api.NewError(http.StatusUnauthorized, "unauthorized")
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("invalid signing method: %v", token.Header["alg"])
			return nil, api.NewError(http.StatusUnauthorized, "unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT: ", err)
		return nil, api.NewError(http.StatusUnauthorized, "unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, api.NewError(http.StatusUnauthorized, "unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println(claims)
		return nil, api.NewError(http.StatusUnauthorized, "unauthorized")
	}
	return claims, nil
}
