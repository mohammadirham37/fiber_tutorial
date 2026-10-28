package resources

import (
	"time"

	"fiberapp/app/models"
)

// File ini setara dengan API Resource di Laravel (php artisan make:resource).
// Tugasnya: menentukan bentuk JSON yang dikirim ke client, supaya kita
// tidak asal mengirim isi model (misalnya email penulis).

type UserResource struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type AuthorResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PostResource struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Author    AuthorResource `json:"author"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// User mengubah model User menjadi JSON untuk pemilik akun (email ikut tampil).
func User(u *models.User) UserResource {
	return UserResource{ID: u.ID, Name: u.Name, Email: u.Email}
}

// Post mengubah satu model Post menjadi PostResource.
func Post(p *models.Post) PostResource {
	return PostResource{
		ID:        p.ID,
		Title:     p.Title,
		Body:      p.Body,
		Author:    AuthorResource{ID: p.User.ID, Name: p.User.Name},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// Posts mengubah banyak Post sekaligus. Setara: PostResource::collection($posts)
func Posts(posts []models.Post) []PostResource {
	result := make([]PostResource, 0, len(posts))
	for i := range posts {
		result = append(result, Post(&posts[i]))
	}
	return result
}
