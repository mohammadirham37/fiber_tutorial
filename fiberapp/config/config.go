package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load membaca file .env ke environment variable.
// Mirip Laravel yang otomatis membaca .env saat aplikasi start.
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, memakai environment sistem")
	}
}

// Env mengambil nilai environment, atau fallback jika kosong.
// Setara dengan helper env('KEY', 'default') di Laravel.
func Env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
