package controllers

import "github.com/gofiber/fiber/v2"

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

// Index menampilkan halaman beranda.
func (h *HomeController) Index(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{
		"title": "Beranda",
	})
}
