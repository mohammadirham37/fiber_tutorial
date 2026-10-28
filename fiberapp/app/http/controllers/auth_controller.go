package controllers

import (
	"fiberapp/app/http/requests"
	"fiberapp/app/models"
	"fiberapp/app/support"
	"fiberapp/database"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// ShowRegister: GET /register
func (a *AuthController) ShowRegister(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"title":  "Daftar",
		"old":    requests.RegisterRequest{},
		"errors": map[string]string{},
	})
}

// Register: POST /register
func (a *AuthController) Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	errs := support.Validate(req)

	// Aturan "unique:users,email" kita cek manual.
	var count int64
	database.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		errs["email"] = "Email sudah terdaftar."
	}

	if len(errs) > 0 {
		req.Password, req.PasswordConfirmation = "", "" // jangan kirim balik password
		return c.Status(fiber.StatusUnprocessableEntity).Render("auth/register", fiber.Map{
			"title":  "Daftar",
			"old":    req,
			"errors": errs,
		})
	}

	// Hash::make($password)
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{Name: req.Name, Email: req.Email, Password: string(hash)}
	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}

	login(c, &user)
	support.Flash(c, "success", "Selamat datang, "+user.Name+"!")
	return c.Redirect("/posts")
}

// ShowLogin: GET /login
func (a *AuthController) ShowLogin(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"title":  "Login",
		"old":    requests.LoginRequest{},
		"errors": map[string]string{},
	})
}

// Login: POST /login
func (a *AuthController) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	errs := support.Validate(req)
	if len(errs) == 0 {
		var user models.User
		err := database.DB.Where("email = ?", req.Email).First(&user).Error

		// Hash::check($password, $user->password)
		if err == nil && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) == nil {
			login(c, &user)
			support.Flash(c, "success", "Berhasil login. Halo, "+user.Name+"!")
			return c.Redirect("/posts")
		}
		errs["email"] = "Email atau password salah."
	}

	req.Password = ""
	return c.Status(fiber.StatusUnprocessableEntity).Render("auth/login", fiber.Map{
		"title":  "Login",
		"old":    req,
		"errors": errs,
	})
}

// Logout: POST /logout
func (a *AuthController) Logout(c *fiber.Ctx) error {
	sess := support.Session(c)
	sess.Delete("user_id")
	sess.Regenerate() // ganti ID session agar session lama tidak bisa dipakai lagi

	support.Flash(c, "success", "Anda sudah logout.")
	return c.Redirect("/login")
}

// login menyimpan ID user ke session. Setara: Auth::login($user)
func login(c *fiber.Ctx, user *models.User) {
	sess := support.Session(c)
	sess.Regenerate() // cegah serangan session fixation
	sess.Set("user_id", user.ID)
}
