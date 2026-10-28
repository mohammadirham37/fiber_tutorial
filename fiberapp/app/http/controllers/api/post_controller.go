package api

import (
	"fmt"

	"fiberapp/app/http/requests"
	"fiberapp/app/http/resources"
	"fiberapp/app/models"
	"fiberapp/database"

	"github.com/gofiber/fiber/v2"
)

// PostController adalah API resource controller untuk Post.
// Setara: php artisan make:controller Api/PostController --api
type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

// Index: GET /api/posts?page=1&per_page=10 — daftar post dengan pagination.
func (p *PostController) Index(c *fiber.Ctx) error {
	page := max(c.QueryInt("page", 1), 1)
	perPage := min(max(c.QueryInt("per_page", 10), 1), 100)

	var total int64
	database.DB.Model(&models.Post{}).Count(&total)

	var posts []models.Post
	database.DB.Preload("User").
		Order("created_at desc").
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(&posts)

	lastPage := (int(total) + perPage - 1) / perPage

	return c.JSON(fiber.Map{
		"data": resources.Posts(posts),
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     perPage,
			"total":        total,
			"last_page":    lastPage,
		},
	})
}

// Show: GET /api/posts/:id — detail satu post.
func (p *PostController) Show(c *fiber.Ctx) error {
	post, err := findPost(c)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": resources.Post(post)})
}

// Store: POST /api/posts — buat post baru (butuh token).
func (p *PostController) Store(c *fiber.Ctx) error {
	var req requests.PostRequest
	if ok, err := parseAndValidate(c, &req); !ok {
		return err
	}

	user := currentUser(c)
	post := models.Post{Title: req.Title, Body: req.Body, UserID: user.ID, User: *user}
	if err := database.DB.Omit("User").Create(&post).Error; err != nil {
		return err
	}

	// 201 Created + header Location berisi alamat data baru.
	c.Location(fmt.Sprintf("/api/posts/%d", post.ID))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": resources.Post(&post)})
}

// Update: PUT/PATCH /api/posts/:id — ubah post milik sendiri.
func (p *PostController) Update(c *fiber.Ctx) error {
	post, err := findOwnedPost(c)
	if err != nil {
		return err
	}

	var req requests.PostRequest
	if ok, err := parseAndValidate(c, &req); !ok {
		return err
	}

	post.Title = req.Title
	post.Body = req.Body
	if err := database.DB.Omit("User").Save(post).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{"data": resources.Post(post)})
}

// Destroy: DELETE /api/posts/:id — hapus post milik sendiri.
func (p *PostController) Destroy(c *fiber.Ctx) error {
	post, err := findOwnedPost(c)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(post).Error; err != nil {
		return err
	}

	// 204 No Content: berhasil, tapi tidak ada isi yang dikirim balik.
	return c.SendStatus(fiber.StatusNoContent)
}

// findPost mencari post berdasarkan :id, atau error 404.
func findPost(c *fiber.Ctx) (*models.Post, error) {
	id, err := c.ParamsInt("id")
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Post tidak ditemukan.")
	}

	var post models.Post
	if err := database.DB.Preload("User").First(&post, id).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Post tidak ditemukan.")
	}
	return &post, nil
}

// findOwnedPost seperti findPost, tapi hanya pemilik post yang boleh (403).
func findOwnedPost(c *fiber.Ctx) (*models.Post, error) {
	post, err := findPost(c)
	if err != nil {
		return nil, err
	}

	if post.UserID != currentUser(c).ID {
		return nil, fiber.NewError(fiber.StatusForbidden, "Anda tidak berhak mengubah post ini.")
	}
	return post, nil
}
