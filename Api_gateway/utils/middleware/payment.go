package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func PaymentAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("UserToken")

		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized, no token present",
			})
		}
		token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("TokenSecret")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}
		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "claims not ok"})
		}
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.UserEmail)

		return c.Next()
	}
}
