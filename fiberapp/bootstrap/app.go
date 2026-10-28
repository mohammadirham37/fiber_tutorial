package bootstrap

import (
	"errors"
	"strings"
	"time"

	"fiberapp/app/http/middleware"
	"fiberapp/config"
	"fiberapp/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

// NewApp merakit aplikasi: view engine, middleware global, dan route.
// Mirip bootstrap/app.php di Laravel.
func NewApp() *fiber.App {
	// 1. View engine (pengganti Blade)
	engine := html.New("./resources/views", ".html")
	engine.Reload(config.Env("APP_DEBUG", "false") == "true") // auto reload view saat development
	engine.AddFunc("limit", limit)

	// 2. Buat aplikasi Fiber
	app := fiber.New(fiber.Config{
		AppName:           config.Env("APP_NAME", "Fiber App"),
		Views:             engine,
		ViewsLayout:       "layouts/main", // layout default semua halaman
		PassLocalsToViews: true,           // c.Locals(...) otomatis bisa dibaca di view
		ErrorHandler:      errorHandler,
	})

	// 3. Middleware global (urutan penting!)
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.MethodOverride())
	app.Use(middleware.StartSession())
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "form:_token", // nama input token, sama seperti Laravel
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		Expiration:     2 * time.Hour,
		ContextKey:     "csrf", // token bisa dibaca di view lewat {{.csrf}}
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api") // API tidak memakai CSRF
		},
	}))
	app.Use(middleware.ShareViewData())

	// 4. File statis dari folder public (CSS, JS, gambar).
	// Catatan: di Fiber v2, app.Static("/") yang diletakkan SEBELUM app.Use()
	// membuat middleware berjalan dua kali. Jadi letakkan setelah middleware.
	app.Static("/", "./public")

	// 5. Daftarkan route
	routes.Web(app)
	routes.Api(app)

	return app
}

// errorHandler menampilkan halaman error (404, 403, 500, ...).
// Untuk /api dikembalikan dalam bentuk JSON.
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	if strings.HasPrefix(c.Path(), "/api") {
		return c.Status(code).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(code).Render("errors/error", fiber.Map{
		"title":   "Terjadi Kesalahan",
		"code":    code,
		"message": err.Error(),
	})
}

// limit memotong teks panjang. Setara Str::limit($text, 100)
func limit(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return string(runes[:max]) + "..."
}
