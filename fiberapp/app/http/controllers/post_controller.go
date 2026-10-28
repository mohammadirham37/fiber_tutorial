package controllers

import (
	"fiberapp/app/http/requests"
	"fiberapp/app/models"
	"fiberapp/app/support"
	"fiberapp/database"

	"github.com/gofiber/fiber/v2"
)

// PostController adalah resource controller untuk Post.
// Setara: php artisan make:controller PostController --resource
type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

// Index: GET /posts — daftar semua post.
func (p *PostController) Index(c *fiber.Ctx) error {
	var posts []models.Post
	database.DB.Preload("User").Order("created_at desc").Find(&posts)

	return c.Render("posts/index", fiber.Map{
		"title": "Semua Post",
		"posts": posts,
	})
}

// Create: GET /posts/create — form tambah post.
func (p *PostController) Create(c *fiber.Ctx) error {
	return c.Render("posts/create", fiber.Map{
		"title":  "Tulis Post Baru",
		"old":    requests.PostRequest{},
		"errors": map[string]string{},
	})
}

// Store: POST /posts — simpan post baru.
func (p *PostController) Store(c *fiber.Ctx) error {
	var req requests.PostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if errs := support.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).Render("posts/create", fiber.Map{
			"title":  "Tulis Post Baru",
			"old":    req,
			"errors": errs,
		})
	}

	user := c.Locals("authUser").(*models.User)
	post := models.Post{
		Title:  req.Title,
		Body:   req.Body,
		UserID: user.ID,
	}
	if err := database.DB.Create(&post).Error; err != nil {
		return err
	}

	support.Flash(c, "success", "Post berhasil dibuat!")
	return c.Redirect("/posts")
}

// Show: GET /posts/:id — detail satu post.
func (p *PostController) Show(c *fiber.Ctx) error {
	post, err := findPost(c)
	if err != nil {
		return err
	}

	return c.Render("posts/show", fiber.Map{
		"title": post.Title,
		"post":  post,
	})
}

// Edit: GET /posts/:id/edit — form edit post.
func (p *PostController) Edit(c *fiber.Ctx) error {
	post, err := findOwnedPost(c)
	if err != nil {
		return err
	}

	return c.Render("posts/edit", fiber.Map{
		"title":  "Edit Post",
		"post":   post,
		"old":    requests.PostRequest{Title: post.Title, Body: post.Body},
		"errors": map[string]string{},
	})
}

// Update: PUT /posts/:id — simpan perubahan post.
func (p *PostController) Update(c *fiber.Ctx) error {
	post, err := findOwnedPost(c)
	if err != nil {
		return err
	}

	var req requests.PostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if errs := support.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).Render("posts/edit", fiber.Map{
			"title":  "Edit Post",
			"post":   post,
			"old":    req,
			"errors": errs,
		})
	}

	post.Title = req.Title
	post.Body = req.Body
	if err := database.DB.Save(post).Error; err != nil {
		return err
	}

	support.Flash(c, "success", "Post berhasil diperbarui!")
	return c.Redirect("/posts/" + c.Params("id"))
}

// Destroy: DELETE /posts/:id — hapus post.
func (p *PostController) Destroy(c *fiber.Ctx) error {
	post, err := findOwnedPost(c)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(post).Error; err != nil {
		return err
	}

	support.Flash(c, "success", "Post berhasil dihapus!")
	return c.Redirect("/posts")
}

// findPost mencari post berdasarkan :id, atau error 404.
// Setara: Post::findOrFail($id)
func findPost(c *fiber.Ctx) (*models.Post, error) {
	id, err := c.ParamsInt("id")
	if err != nil {
		return nil, fiber.ErrNotFound
	}

	var post models.Post
	if err := database.DB.Preload("User").First(&post, id).Error; err != nil {
		return nil, fiber.ErrNotFound
	}
	return &post, nil
}

// findOwnedPost sama seperti findPost, tapi hanya pemilik post yang boleh.
// Setara: $this->authorize('update', $post) dengan Policy.
func findOwnedPost(c *fiber.Ctx) (*models.Post, error) {
	post, err := findPost(c)
	if err != nil {
		return nil, err
	}

	user := c.Locals("authUser").(*models.User)
	if post.UserID != user.ID {
		return nil, fiber.NewError(fiber.StatusForbidden, "Anda tidak berhak mengubah post ini.")
	}
	return post, nil
}
