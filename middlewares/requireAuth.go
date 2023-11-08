package middlewares

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmcerez0/fiber-demo/models"
	"github.com/jmcerez0/fiber-demo/utils"
)

func RequireAuth(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")

	if tokenString == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user models.User
		utils.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("user", user)

		c.Next()
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return err
}
