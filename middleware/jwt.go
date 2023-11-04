package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/GrigoryNazarov96/hotel-reservation.git/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(s db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Bearer"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}
		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired, please login again")
		}
		userId := claims["id"].(string)
		user, err := s.GetUserByID(c.Context(), userId)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("invalid signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT: ", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println(claims)
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}
