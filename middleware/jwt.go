package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["Bearer"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if err := parseToken(token[0]); err != nil {
		return err
	}
	return nil
}

func parseToken(tokenStr string) error {
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
		return fmt.Errorf("unauthorized")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		return nil
	}
	return fmt.Errorf("unauthorized")
}
