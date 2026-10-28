package middleware

import (
	"strings"

	"fiberapp/app/support"

	"github.com/gofiber/fiber/v2"
)

// ApiAuth melindungi route API dengan token Bearer.
// Setara: middleware('auth:sanctum') di Laravel.
//
// Client harus mengirim header:  Authorization: Bearer <token>
func ApiAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		plain, found := strings.CutPrefix(header, "Bearer ")
		if !found || plain == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Token tidak ditemukan.")
		}

		token, err := support.FindToken(plain)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Token tidak valid.")
		}

		c.Locals("apiUser", &token.User)
		c.Locals("apiToken", token)
		return c.Next()
	}
}
