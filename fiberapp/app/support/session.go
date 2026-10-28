package support

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Store menyimpan seluruh session (default: di memori).
var Store = session.New(session.Config{
	Expiration:     2 * time.Hour,
	CookieHTTPOnly: true,
	CookieSameSite: "Lax",
})

// Session mengambil session milik request saat ini.
// Session dibuka sekali oleh middleware StartSession lalu dipakai ulang.
func Session(c *fiber.Ctx) *session.Session {
	return c.Locals("session").(*session.Session)
}

// Flash menyimpan pesan yang hanya tampil SATU kali di request berikutnya.
// Setara: session()->flash('success', 'Pesan') di Laravel.
func Flash(c *fiber.Ctx, key, message string) {
	Session(c).Set("flash_"+key, message)
}

// PullFlash mengambil pesan flash lalu langsung menghapusnya.
func PullFlash(c *fiber.Ctx, key string) string {
	sess := Session(c)
	message, _ := sess.Get("flash_" + key).(string)
	sess.Delete("flash_" + key)
	return message
}
