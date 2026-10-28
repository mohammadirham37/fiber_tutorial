# Modul Ajar: Fullstack Go Fiber ala Laravel

Modul belajar berbahasa Indonesia untuk pemula: mulai dari dasar bahasa Go, dasar web (HTTP), sampai membangun aplikasi blog lengkap dengan register, login, CRUD, dan RESTful API memakai [Fiber](https://gofiber.io) dengan struktur folder ala Laravel.

## Isi repository

| Path | Isi |
| --- | --- |
| `modul-fiber-laravel-style.html` | Modul ajar (buka langsung di browser) |
| `fiberapp/` | Kode lengkap aplikasi hasil akhir modul |

## Daftar materi

1. **Persiapan** — mengenal Go & Fiber, instalasi
2. **Dasar Bahasa Go** — variabel, percabangan, perulangan, slice, map, function, error, struct, pointer, interface, package, mini proyek CRUD di terminal
3. **Dasar Web & Fiber** — HTTP, request, response, server Fiber pertama
4. **Studi Kasus Blog** — `.env`, GORM, migration, seeder, routing, controller, view, session, flash, middleware, CSRF, validasi, register/login/logout, CRUD dengan cek pemilik
5. **RESTful API** — konsep REST, token ala Sanctum, CRUD JSON, pagination, API Resource
6. **Lanjutan** — perintah ala Artisan, hot reload, deploy, troubleshooting, latihan

## Menjalankan aplikasi contoh

```bash
cd fiberapp
cp .env.example .env
go run . migrate
go run . db:seed
go run .
```

Buka http://localhost:3000 dan login dengan `admin@example.com` / `password`.

Contoh request API ada di `fiberapp/api.http` (extension VS Code "REST Client").
