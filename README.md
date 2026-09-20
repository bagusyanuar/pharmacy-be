# Go Clean Architecture & DDD Backend Template

A production-ready, batteries-included Golang backend starter template featuring **Domain-Driven Design (DDD)** and **Clean Architecture**, with built-in **AI pair programming skills** (`.agents/`).

---

## 🌟 Key Features

- **Clean Architecture & DDD**: Strict layer separation (`domain`, `usecase`, `infrastructure`, `delivery`).
- **Framework & Performance**: Powered by [Fiber v3](https://gofiber.io/) and [GORM](https://gorm.io/) (PostgreSQL default, multi-driver ready).
- **Authentication**: JWT access token + secure `HttpOnly` refresh token cookie.
- **Built-in Portable AI Skills**: Comes with `.agents/skills/go-module-generator` and architectural rules in `.agents/rules/` for Antigravity / Gemini / Claude Code / Cursor.
- **Robust Cross-Module Composition**: Clean boundary orchestration without cross-module GORM joins/preloads.
- **Structured Logging & Diagnostics**: Uber [Zap](https://github.com/uber-go/zap) + [Lumberjack](https://github.com/natefinch/lumberjack) log rotation.
- **Standardized API Envelope**: Uniform JSON responses (`data`, `meta`, `error`) and validation error details (Laravel-style).
- **Docker & Tooling**: Includes `Dockerfile`, `docker-compose.yml`, and `Makefile`.

---

## 📖 Panduan Penggunaan untuk Project Baru (How to Use for a New Project)

Ada 4 opsi mudah untuk memulai project baru menggunakan template ini:

### Opsi 1: Menggunakan `gonew` (Cara Resmi Go - Paling Cepat & Modern) ⚡
Go memiliki tool resmi bernama [gonew](https://pkg.go.dev/golang.org/x/tools/cmd/gonew) untuk scaffolding project dari template:

1. Install `gonew` (cukup sekali di komputermu):
   ```bash
   go install golang.org/x/tools/cmd/gonew@latest
   ```

2. Jalankan perintah pembuatan project baru dari folder mana saja di terminal:
   ```bash
   # Format: gonew <template-repo> <target-module-name> [target-directory]
   gonew github.com/bagusyanuar/go-be-template github.com/username/my-new-service my-new-service
   ```
   *Secara otomatis `gonew` akan mengunduh template, mengganti semua module & import paths ke module barumu, dan menyiapkannya di folder `my-new-service`.*

3. Masuk ke folder baru dan mulai development:
   ```bash
   cd my-new-service
   cp .env.example .env
   make run
   ```

---

### Opsi 2: Menggunakan GitHub Template (Paling Mudah via Browser)
1. Buka repositori template ini di GitHub, lalu klik tombol hijau **"Use this template"** > **"Create a new repository"**.
2. Beri nama repositori baru kamu (contoh: `payment-service`) dan clone ke komputer:
   ```bash
   git clone https://github.com/username/payment-service.git
   cd payment-service
   ```
3. Jalankan command inisialisasi otomatis:
   ```bash
   make init MODULE=github.com/username/payment-service
   ```
   *(Perintah ini akan mengganti seluruh import path di semua file Go, membuat file `.env`, dan menjalankan `go mod tidy`)*.

---

### Opsi 3: Clone Manual / Copy Folder
Jika belum di-push ke GitHub atau ingin langsung dari lokal:
1. Salin/clone template ke folder project baru:
   ```bash
   # Contoh clone ke folder project baru di luar template
   git clone https://github.com/bagusyanuar/go-be-template my-new-service
   cd my-new-service
   
   # Reset git history agar menjadi repo baru
   rm -rf .git
   git init
   ```
2. Jalankan perintah `make init`:
   ```bash
   make init MODULE=github.com/username/my-new-service
   ```

---

### Opsi 4: Langsung Menggunakan AI di IDE (Antigravity / Gemini)
Karena template ini sudah membawa **AI Rules & Skills** di dalam folder `.agents/`, kamu cukup:
1. Buka folder project hasil clone di IDE.
2. Buka chat AI, lalu ketik:
   > *"Bro, tolong inisialisasi project ini untuk module `github.com/username/order-service`"*
3. AI akan otomatis mengeksekusi script inisialisasi dan menyiapkan environment-nya untukmu!

---

## 🚀 Menjalankan Project

### 1. Menggunakan Docker (Rekomendasi Cepat)
Menjalankan database PostgreSQL dan API server sekaligus:
```bash
make docker-up
```
Aplikasi langsung aktif di `http://localhost:3000`.

### 2. Menjalankan secara Lokal
Pastikan kamu sudah memiliki PostgreSQL lokal yang aktif:
```bash
# 1. Pastikan .env sudah sesuai dengan koneksi database lokal kamu
cp .env.example .env

# 2. Jalankan server
make run
```

### 3. Menjalankan Unit Tests
```bash
make test
```

---

## 🤖 Menambah Modul Baru (Scaffolding with AI)

Template ini memiliki **AI Skill Generator Modul DDD** bawaan. Kamu tidak perlu mengetik boilerplate 4-layer secara manual.

Cukup chat ke AI pair programming:
> *"Tolong buatkan module baru `product` dengan fields: `name` (string), `sku` (string, unique), `price` (int64), `stock` (int), dan sediakan usecase CRUD lengkap."*

**AI akan otomatis:**
1. Membuat `domain/product.go` & `domain/errors.go` (Entity, Repository Interface, Sentinel Errors).
2. Membuat `infrastructure/product_repository.go` (GORM query dengan context awareness & paging).
3. Membuat `usecase/product_usecase.go` (Business logic flow & DTOs).
4. Membuat `delivery/product_handler.go` & `delivery/product_dto.go` (Validasi input & JSON envelope response).
5. Mendaftarkan dependency injection-nya langsung ke `internal/bootstrap/bootstrap.go`.

---

## 📁 Struktur Direktori

```text
├── .agents/
│   ├── rules/                      # Aturan arsitektur, validasi, naming, query, dan error
│   └── skills/
│       └── go-module-generator/    # AI skill untuk generate modul DDD otomatis
├── AGENTS.md                       # Entrypoint instruksi untuk AI coding assistants
├── cmd/
│   └── api/
│       └── main.go                 # Process lifecycle & graceful shutdown
├── internal/
│   ├── bootstrap/                  # DI wiring (Repo -> Usecase -> Handler) & router
│   ├── config/                     # Environment configuration & PostgreSQL pool
│   ├── middleware/                 # Fiber middlewares (Auth, Logger, CORS)
│   └── modules/                    # Bounded Contexts (DDD)
│       ├── auth/                   # Modul autentikasi (Register, Login, Refresh, Me)
│       └── user/                   # Modul user (Domain entity, Repo, Usecase, Handler)
├── pkg/                            # Shared utilities (jwt, logger, password, response, validator)
├── scripts/
│   └── init-project.sh             # Skrip inisialisasi module untuk project baru
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

---

## 📜 Default API Endpoints

| Method | Endpoint | Description | Auth Required |
|---|---|---|---|
| `GET` | `/health` | Service health status | No |
| `POST` | `/api/v1/auth/register` | Register new account | No |
| `POST` | `/api/v1/auth/login` | Login and receive tokens | No |
| `POST` | `/api/v1/auth/refresh` | Refresh access token via cookie | No (Cookie) |
| `POST` | `/api/v1/auth/logout` | Clear refresh token cookie | No |
| `GET` | `/api/v1/auth/me` | Current user profile | Yes (Bearer) |
| `GET` | `/api/v1/users` | List users (paginated) | Yes (Bearer) |
| `GET` | `/api/v1/users/:id` | Get user by ID | Yes (Bearer) |
