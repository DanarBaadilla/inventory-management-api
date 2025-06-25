# 🗃️ Inventory Management API

RESTful API untuk mengelola sistem inventaris seperti kategori produk, produk, pergerakan stok, dan manajemen user. Dibangun dengan bahasa Go menggunakan framework Fiber serta mengikuti prinsip arsitektur clean code.

-----

## 🔧 Teknologi yang Digunakan

- [Go Fiber](https://gofiber.io/) — Web Framework
- [GORM](https://gorm.io/) — ORM untuk koneksi MySQL
- [MySQL](https://www.mysql.com/) — Database Relasional
- [JWT](https://jwt.io/) — Autentikasi
- [Swagger (Swaggo)](https://github.com/swaggo/swag) — Dokumentasi API
- [Validator.v10](https://pkg.go.dev/github.com/go-playground/validator/v10) — Validasi input
- Clean Architecture

-----

## 📁 Struktur Folder

```
inventory-management-api/
├── app/ # Inisialisasi Fiber app
├── config/ # Koneksi database dan konfigurasi
├── controller/ # HTTP handler
├── docs/ # Dokumentasi Swagger
├── helper/ # Fungsi utilitas/helper
├── middleware/ # JWT, CORS, dll
├── model/
│ ├── domain/ # Model untuk database (GORM)
│ ├── web/ # Model untuk request/response
├── repository/ # Akses data ke database
├── route/ # Routing endpoint
├── service/ # Logika bisnis
├── main.go # Entry point
```

-----

## 🔒 Autentikasi

Gunakan token JWT pada header `Authorization`:

```http
Authorization: Bearer <token>
```

-----

## 📄 Dokumentasi Swagger

Akses dokumentasi di:
👉 **http://localhost:3000/swagger/index.html**

-----

## 📝 Lisensi

MIT License © 2025 Danar Rafiardi