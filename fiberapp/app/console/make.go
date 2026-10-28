package console

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const controllerStub = `package controllers

import "github.com/gofiber/fiber/v2"

type {{Name}} struct{}

func New{{Name}}() *{{Name}} {
	return &{{Name}}{}
}

func (h *{{Name}}) Index(c *fiber.Ctx) error {
	return c.SendString("Halo dari {{Name}}")
}
`

// MakeController membuat file controller baru dari template (stub).
// Setara: php artisan make:controller NamaController
func MakeController(name string) error {
	if !strings.HasSuffix(name, "Controller") {
		name += "Controller"
	}

	path := filepath.Join("app", "http", "controllers", toSnake(name)+".go")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("controller sudah ada: %s", path)
	}

	content := strings.ReplaceAll(controllerStub, "{{Name}}", name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	fmt.Println("Controller dibuat:", path)
	return nil
}

// toSnake mengubah "ProductController" menjadi "product_controller".
func toSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			r = unicode.ToLower(r)
		}
		b.WriteRune(r)
	}
	return b.String()
}
