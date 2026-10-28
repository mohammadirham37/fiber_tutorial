package middleware

import (
	"fiberapp/app/support"

	"github.com/gofiber/fiber/v2"
)

// Auth hanya mengizinkan user yang sudah login. Setara middleware('auth').
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("authUser") == nil {
			support.Flash(c, "error", "Silakan login terlebih dahulu.")
			return c.Redirect("/login")
		}
		return c.Next()
	}
}

// Guest hanya mengizinkan user yang BELUM login. Setara middleware('guest').
func Guest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("authUser") != nil {
			return c.Redirect("/posts")
		}
		return c.Next()
	}
}
