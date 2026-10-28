package support

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	// Pakai nama dari tag `form` (misal "title") untuk pesan error,
	// bukan nama field Go (misal "Title").
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("form")
	})
}

// Validate memeriksa struct dan mengembalikan map: nama_field => pesan error.
// Map kosong berarti data valid.
func Validate(data any) map[string]string {
	result := map[string]string{}

	var validationErrors validator.ValidationErrors
	if err := validate.Struct(data); errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			result[e.Field()] = message(e)
		}
	}
	return result
}

// message menerjemahkan aturan validasi menjadi pesan berbahasa Indonesia.
func message(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Kolom %s wajib diisi.", e.Field())
	case "email":
		return fmt.Sprintf("Kolom %s harus berupa alamat email yang valid.", e.Field())
	case "min":
		return fmt.Sprintf("Kolom %s minimal %s karakter.", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("Kolom %s maksimal %s karakter.", e.Field(), e.Param())
	case "eqfield":
		return fmt.Sprintf("Kolom %s harus sama dengan kolom %s.", e.Field(), strings.ToLower(e.Param()))
	}
	return fmt.Sprintf("Kolom %s tidak valid.", e.Field())
}
