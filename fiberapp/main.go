package main

import (
	"fmt"
	"log"
	"os"

	"fiberapp/app/console"
	"fiberapp/bootstrap"
	"fiberapp/config"
	"fiberapp/database"
)

// main bekerja seperti "php artisan":
//
//	go run . serve            -> jalankan server
//	go run . migrate          -> buat tabel
//	go run . migrate:fresh    -> hapus semua tabel lalu buat ulang
//	go run . db:seed          -> isi data awal
//	go run . make:controller Product
func main() {
	config.Load()

	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "serve":
		database.Connect()
		app := bootstrap.NewApp()
		log.Fatal(app.Listen(":" + config.Env("APP_PORT", "3000")))

	case "migrate":
		database.Connect()
		database.Migrate()

	case "migrate:fresh":
		database.Connect()
		database.Fresh()

	case "db:seed":
		database.Connect()
		database.Seed()

	case "make:controller":
		if len(os.Args) < 3 {
			log.Fatal("Contoh: go run . make:controller Product")
		}
		if err := console.MakeController(os.Args[2]); err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("Perintah tersedia: serve, migrate, migrate:fresh, db:seed, make:controller")
	}
}
