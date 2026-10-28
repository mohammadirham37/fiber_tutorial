package routes

import (
	"fiberapp/app/http/controllers"
	"fiberapp/app/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// Web mendaftarkan semua route halaman web. Setara routes/web.php
func Web(app *fiber.App) {
	home := controllers.NewHomeController()
	auth := controllers.NewAuthController()
	post := controllers.NewPostController()

	app.Get("/", home.Index).Name("home")

	// Auth (hanya untuk tamu / belum login)
	app.Get("/register", middleware.Guest(), auth.ShowRegister).Name("register")
	app.Post("/register", middleware.Guest(), auth.Register)
	app.Get("/login", middleware.Guest(), auth.ShowLogin).Name("login")
	app.Post("/login", middleware.Guest(), auth.Login)
	app.Post("/logout", middleware.Auth(), auth.Logout).Name("logout")

	// Resource posts. Urutan penting: /posts/create harus di atas /posts/:id
	app.Get("/posts", post.Index).Name("posts.index")
	app.Get("/posts/create", middleware.Auth(), post.Create).Name("posts.create")
	app.Post("/posts", middleware.Auth(), post.Store).Name("posts.store")
	app.Get("/posts/:id", post.Show).Name("posts.show")
	app.Get("/posts/:id/edit", middleware.Auth(), post.Edit).Name("posts.edit")
	app.Put("/posts/:id", middleware.Auth(), post.Update).Name("posts.update")
	app.Delete("/posts/:id", middleware.Auth(), post.Destroy).Name("posts.destroy")
}
