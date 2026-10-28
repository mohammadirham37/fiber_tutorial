package middleware

import (
	"fiberapp/app/models"
	"fiberapp/app/support"
	"fiberapp/database"

	"github.com/gofiber/fiber/v2"
)

// ShareViewData menyiapkan data yang dipakai di SEMUA view:
// user yang sedang login dan pesan flash.
// Mirip View::share() / View Composer di Laravel.
func ShareViewData() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess := support.Session(c)

		if id, ok := sess.Get("user_id").(uint); ok {
			var user models.User
			if err := database.DB.First(&user, id).Error; err == nil {
				c.Locals("authUser", &user)
			}
		}

		c.Locals("success", support.PullFlash(c, "success"))
		c.Locals("error", support.PullFlash(c, "error"))

		return c.Next()
	}
}
