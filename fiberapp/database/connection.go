package database

import (
	"fmt"
	"log"
	"os"

	"fiberapp/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah koneksi database global, mirip facade DB di Laravel.
var DB *gorm.DB

// Connect membuka koneksi sesuai DB_CONNECTION di file .env.
func Connect() {
	var dialector gorm.Dialector

	switch config.Env("DB_CONNECTION", "sqlite") {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Env("DB_USERNAME", "root"),
			config.Env("DB_PASSWORD", ""),
			config.Env("DB_HOST", "127.0.0.1"),
			config.Env("DB_PORT", "3306"),
			config.Env("DB_DATABASE", "fiber_blog"),
		)
		dialector = mysql.Open(dsn)
	default:
		os.MkdirAll("storage", 0755)
		dialector = sqlite.Open(config.Env("DB_DATABASE", "storage/database.sqlite"))
	}

	logLevel := logger.Silent
	if config.Env("APP_DEBUG", "false") == "true" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}

	DB = db
}
