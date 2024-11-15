package helper

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GetRequiredFormValue(c *fiber.Ctx, key string) (string, error) {
	value := c.FormValue(key)
	if value == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, key+" is required")
	}
	return value, nil
}
