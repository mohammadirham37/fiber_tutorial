package api

import (
	"fiberapp/app/models"
	"fiberapp/app/support"

	"github.com/gofiber/fiber/v2"
)

// parseAndValidate membaca body JSON ke dalam req lalu memvalidasinya.
// Jika gagal, langsung mengirim respons 400 atau 422 dan ok bernilai false.
func parseAndValidate(c *fiber.Ctx, req any) (ok bool, err error) {
	if err := c.BodyParser(req); err != nil {
		return false, fiber.NewError(fiber.StatusBadRequest, "Body JSON tidak valid.")
	}

	if errs := support.Validate(req); len(errs) > 0 {
		// Bentuk respons sama seperti Laravel saat validasi gagal.
		return false, c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Data yang dikirim tidak valid.",
			"errors":  errs,
		})
	}
	return true, nil
}

// currentUser mengambil user pemilik token. Setara: $request->user()
func currentUser(c *fiber.Ctx) *models.User {
	return c.Locals("apiUser").(*models.User)
}
