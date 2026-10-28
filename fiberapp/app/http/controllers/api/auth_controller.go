package api

import (
	"fiberapp/app/http/requests"
	"fiberapp/app/http/resources"
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

// Register: POST /api/register — buat akun baru dan langsung beri token.
func (a *AuthController) Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if ok, err := parseAndValidate(c, &req); !ok {
		return err
	}

	var count int64
	database.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Data yang dikirim tidak valid.",
			"errors":  fiber.Map{"email": "Email sudah terdaftar."},
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{Name: req.Name, Email: req.Email, Password: string(hash)}
	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}

	token, err := support.CreateToken(&user, "api")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":  resources.User(&user),
		"token": token,
	})
}

// Login: POST /api/login — tukar email + password dengan token.
func (a *AuthController) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if ok, err := parseAndValidate(c, &req); !ok {
		return err
	}

	var user models.User
	err := database.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Email atau password salah.")
	}

	token, err := support.CreateToken(&user, "api")
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data":  resources.User(&user),
		"token": token,
	})
}

// Me: GET /api/me — data user pemilik token.
func (a *AuthController) Me(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"data": resources.User(currentUser(c))})
}

// Logout: POST /api/logout — hapus token yang sedang dipakai.
func (a *AuthController) Logout(c *fiber.Ctx) error {
	token := c.Locals("apiToken").(*models.PersonalAccessToken)
	if err := database.DB.Delete(token).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Token sudah dihapus. Anda sudah logout."})
}
