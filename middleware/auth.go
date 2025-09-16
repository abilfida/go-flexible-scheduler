package middleware

import (
	"github.com/abilfida/go-flexible-scheduler/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v4"
)

// Protected melindungi rute yang memerlukan otentikasi.
func Protected() fiber.Handler {
	cfg, _ := config.LoadConfig()

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWTSecret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	})
}
