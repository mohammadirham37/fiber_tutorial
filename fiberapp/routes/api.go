package routes

import (
	"fiberapp/app/http/controllers/api"
	"fiberapp/app/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Api mendaftarkan route JSON dengan prefix /api. Setara routes/api.php
func Api(app *fiber.App) {
	auth := api.NewAuthController()
	post := api.NewPostController()

	// cors mengizinkan aplikasi frontend dari domain lain memanggil API ini.
	r := app.Group("/api", cors.New())

	// Publik: tidak perlu token
	r.Post("/register", auth.Register)
	r.Post("/login", auth.Login)
	r.Get("/posts", post.Index)
	r.Get("/posts/:id", post.Show)

	// Terlindungi: wajib kirim header Authorization: Bearer <token>
	// Setara: ->middleware('auth:sanctum') di Laravel
	authed := middleware.ApiAuth()
	r.Get("/me", authed, auth.Me)
	r.Post("/logout", authed, auth.Logout)
	r.Post("/posts", authed, post.Store)
	r.Put("/posts/:id", authed, post.Update)
	r.Patch("/posts/:id", authed, post.Update)
	r.Delete("/posts/:id", authed, post.Destroy)
}
