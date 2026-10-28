package database

import (
	"log"

	"fiberapp/app/models"

	"golang.org/x/crypto/bcrypt"
)

// Seed mengisi data awal. Setara: php artisan db:seed
func Seed() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	admin := models.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: string(hash),
	}
	if err := DB.Where(models.User{Email: admin.Email}).FirstOrCreate(&admin).Error; err != nil {
		log.Fatal("Seeder user gagal: ", err)
	}

	posts := []models.Post{
		{Title: "Belajar Fiber dari Nol", Body: "Fiber adalah web framework Go yang terinspirasi dari Express.js. Cepat, ringan, dan mudah dipelajari.", UserID: admin.ID},
		{Title: "Struktur Folder ala Laravel", Body: "Dengan struktur folder yang rapi, project Go terasa familiar bagi developer Laravel.", UserID: admin.ID},
		{Title: "GORM: Eloquent-nya Golang", Body: "GORM memudahkan kita bekerja dengan database tanpa menulis SQL manual.", UserID: admin.ID},
	}
	for _, post := range posts {
		DB.Where(models.Post{Title: post.Title}).FirstOrCreate(&post)
	}

	log.Println("Seeder selesai. Login: admin@example.com / password")
}
