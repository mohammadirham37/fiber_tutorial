package requests

// Tag `form` dipakai saat data datang dari form HTML.
// Tag `json` dipakai saat data datang dari body JSON (API).
type LoginRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name                 string `form:"name" json:"name" validate:"required,min=3,max=100"`
	Email                string `form:"email" json:"email" validate:"required,email,max=100"`
	Password             string `form:"password" json:"password" validate:"required,min=8"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" validate:"required,eqfield=Password"`
}
