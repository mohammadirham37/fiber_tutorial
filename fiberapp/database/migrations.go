package database

import (
	"log"

	"fiberapp/app/models"
)

// Migrate membuat/memperbarui tabel. Setara: php artisan migrate
func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.PersonalAccessToken{},
	)
	if err != nil {
		log.Fatal("Migrasi gagal: ", err)
	}
	log.Println("Migrasi selesai.")
}

// Fresh menghapus semua tabel lalu migrasi ulang. Setara: php artisan migrate:fresh
func Fresh() {
	if err := DB.Migrator().DropTable(&models.PersonalAccessToken{}, &models.Post{}, &models.User{}); err != nil {
		log.Fatal("Gagal menghapus tabel: ", err)
	}
	log.Println("Semua tabel dihapus.")
	Migrate()
}
