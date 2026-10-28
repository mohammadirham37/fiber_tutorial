package middleware

import (
	"fiberapp/app/support"

	"github.com/gofiber/fiber/v2"
)

// StartSession membuka session di awal request dan menyimpannya
// otomatis di akhir request (seperti middleware StartSession Laravel).
func StartSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := support.Store.Get(c)
		if err != nil {
			return err
		}
		c.Locals("session", sess)

		err = c.Next()

		if saveErr := sess.Save(); saveErr != nil && err == nil {
			return saveErr
		}
		return err
	}
}
