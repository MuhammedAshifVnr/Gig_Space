package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID             uint
	UserEmail          string
	Role               string
	SubscriptionExpiry int64
	jwt.StandardClaims
}

func Auth(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies(role)

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

		// if role == "UserToken" && time.Now().Unix() > claims.SubscriptionExpiry {
		// 	return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
		// 		"message": "Your subscription is expired or inactive. Please renew to access this service.",
		// 	})
		// }
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.UserEmail)

		return c.Next()
	}
}

func AuthChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := strings.TrimPrefix(c.Get("UserToken"), "Bearer ")
		fmt.Println("==")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("TokenSecret")), nil
		})

		if err != nil || !token.Valid {
			log.Println("MW: User not Authorized")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}
		Claims, ok := token.Claims.(*Claims)
		if !ok {
			c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "claims not ok"})
		}

		c.Locals("userID", Claims.UserID)
		c.Locals("User_role", Claims.Role)
		fmt.Println("userID", Claims.UserID)

		return c.Next()
	}
}

func AuthClient(token string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies(token)

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

		// if role == "UserToken" && time.Now().Unix() > claims.SubscriptionExpiry {
		// 	return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
		// 		"message": "Your subscription is expired or inactive. Please renew to access this service.",
		// 	})
		// }
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.UserEmail)

		return c.Next()
	}
}
