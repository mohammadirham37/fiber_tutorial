package requests

// PostRequest berisi data form post beserta aturan validasinya.
// Mirip class Form Request di Laravel (php artisan make:request).
type PostRequest struct {
	Title string `form:"title" json:"title" validate:"required,min=3,max=200"`
	Body  string `form:"body" json:"body" validate:"required,min=10"`
}
