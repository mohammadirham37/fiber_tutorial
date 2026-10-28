package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MethodOverride membuat form HTML bisa mengirim PUT/PATCH/DELETE
// lewat input tersembunyi _method. Setara dengan @method('PUT') di Blade.
//
// PENTING: daftarkan middleware ini SEBELUM semua route.
func MethodOverride() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodPost {
			method := strings.ToUpper(c.FormValue("_method"))
			switch method {
			case fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete:
				c.Method(method)
			}
		}
		return c.Next()
	}
}
