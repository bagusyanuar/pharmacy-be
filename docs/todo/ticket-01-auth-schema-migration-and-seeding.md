# 📋 Todo List: Tiket #1 — Migrasi Skema Database Relasional, Indeks & Seeding

* **Tiket Acuan:** [`bagusyanuar/pharmacy-docs#1`](https://github.com/bagusyanuar/pharmacy-docs/issues/1)
* **Modul:** `[BE] Auth & Sesi (Part 1/5)`
* **Target Branch:** `feat/docs-1-auth-schema-migration`
* **Status:** 🟢 Selesai (Siap Buka PR)

---

## 🎯 Deliverable Checklist

### Fase 1: Setup Tooling Migrasi (`golang-migrate`)
- [x] Tambahkan dependensi migrasi ke `go.mod`:
  - `github.com/golang-migrate/migrate/v4`
  - `github.com/golang-migrate/migrate/v4/database/postgres`
  - `github.com/golang-migrate/migrate/v4/source/file`
- [x] Buat CLI runner migrasi di `cmd/migrate/main.go`.
- [x] Tambahkan target di `Makefile`:
  - `make migrate-up` ➡️ Menjalankan seluruh migrasi yang belum terpasang.
  - `make migrate-down` ➡️ Melakukan rollback 1 tahap migrasi.
  - `make migrate-create NAME=<nama>` ➡️ Membuat pasangan file `.up.sql` dan `.down.sql`.
  - `make migrate-version` ➡️ Memeriksa versi migrasi saat ini.
  - `make migrate-force VERSION=<v>` ➡️ Paksa set versi migrasi.

---

### Fase 2: Pembuatan File Migrasi Granular (DDL PostgreSQL 15+)
Urutan migrasi disusun berdasarkan dependensi *Foreign Key*:
- [x] `database/migrations/000001_create_branches_table.up.sql / .down.sql`
  - Kolom: `id`, `branch_code`, `branch_name`, `branch_type`, `sia_number`, `sia_expired_date`, `phone`, `email`, `address`, `city`, `province`, `postal_code`, `is_active`.
  - 5 Kolom Universal Audit Trail.
- [x] `database/migrations/000002_create_roles_table.up.sql / .down.sql`
  - Kolom: `id`, `role_code`, `role_name`, `description`, `can_supervisor_override`.
  - 5 Kolom Universal Audit Trail.
- [x] `database/migrations/000003_create_users_table.up.sql / .down.sql`
  - Kolom: `id`, `email`, `password_hash`, `pin_hash`, `barcode_card`, `role_id`, `is_active`.
  - FK `role_id REFERENCES roles(id) ON DELETE RESTRICT`.
  - 5 Kolom Universal Audit Trail.
- [x] `database/migrations/000004_create_staff_profiles_table.up.sql / .down.sql`
  - Kolom: `id`, `user_id`, `employee_code`, `nik`, `full_name`, `gender`, `phone`, `email`, `address`, `hire_date`, `job_position`, `is_active`.
  - FK `user_id REFERENCES users(id) ON DELETE SET NULL` (Relasi 1-to-1 opsional).
  - 5 Kolom Universal Audit Trail.
- [x] `database/migrations/000005_create_user_branches_table.up.sql / .down.sql`
  - Kolom: `id`, `user_id`, `branch_id`, `is_default`, `can_operate_pos`.
  - FK `user_id REFERENCES users(id) ON DELETE RESTRICT`.
  - FK `branch_id REFERENCES branches(id) ON DELETE RESTRICT`.
  - Composite Unique: `uq_user_branches_assignment (user_id, branch_id)`.
  - 5 Kolom Universal Audit Trail.
- [x] `database/migrations/000006_create_pharmacist_profiles_table.up.sql / .down.sql`
  - Kolom: `id`, `staff_id`, `license_type`, `license_number`, `license_expired_date`, `is_apa`, `assigned_branch_id`.
  - FK `staff_id REFERENCES staff_profiles(id) ON DELETE RESTRICT`.
  - FK `assigned_branch_id REFERENCES branches(id) ON DELETE RESTRICT`.
  - 5 Kolom Universal Audit Trail.
- [x] `database/migrations/000007_create_supervisor_override_logs_table.up.sql / .down.sql`
  - Kolom: `id`, `branch_id`, `cashier_user_id`, `supervisor_user_id`, `action_type`, `reason_category`, `reason_notes`, `target_entity_type`, `target_entity_id`.
  - FK `branch_id REFERENCES branches(id) ON DELETE RESTRICT`.
  - FK `cashier_user_id REFERENCES users(id) ON DELETE RESTRICT`.
  - FK `supervisor_user_id REFERENCES users(id) ON DELETE RESTRICT`.
  - 5 Kolom Universal Audit Trail.

---

### Fase 3: Penegakan Invarian & Indeks Performa Tinggi
- [x] Buat Partial Unique Index (`WHERE deleted_at IS NULL`):
  - `uq_staff_user_id` pada `staff_profiles(user_id)`
  - `uq_staff_employee_code` pada `staff_profiles(employee_code)`
  - `uq_staff_nik` pada `staff_profiles(nik)`
  - `uq_users_barcode_card` pada `users(barcode_card)`
  - `uq_users_email` pada `users(email)`
- [x] Buat Composite Performance Index:
  - `idx_override_logs_branch_date` pada `supervisor_override_logs(branch_id, created_at DESC)`

---

### Fase 4: Penyelarasan Domain (Refactor Scaffold Lama Sesuai DDD Bounded Context)
- [x] Refactor `internal/modules/user/domain/user.go`:
  - Hapus field `username` (sesuai SSOT baru yang login via email/barcode).
  - Tambahkan `pin_hash`, `barcode_card`, `role_id`, `created_by`.
  - Hapus import `gorm.DeletedAt` di layer domain (ganti `*time.Time` murni).
- [x] Definisikan domain entities murni sesuai batas *Bounded Context*:
  - **IAM Context**: `Role`, `UserBranch` (`internal/modules/user/domain/`)
  - **Branch Core Context**: `Branch` (`internal/modules/branch/domain/branch.go`)
  - **HR / Staff Context**: `StaffProfile`, `PharmacistProfile` (`internal/modules/staff/domain/`)
  - **Auth & Audit Context**: `SupervisorOverrideLog` (`internal/modules/auth/domain/supervisor_override_log.go`)

---

### Fase 5: Data Seeding Awal Idempotent
- [x] Buat runner seeder di `database/seeders/` yang dapat dijalankan via `make seed` (`cmd/seed/main.go`).
- [x] Seed 6 Role Baku Sistem (`SUPER_ADMIN`, `APOTEKER`, `ASISTEN_APOTEKER`, `KASIR`, `GUDANG`, `FINANCE_OWNER`).
- [x] Seed 1 Cabang Apotek Default di `branches`.
- [x] Seed 1 Akun Default `SUPER_ADMIN` di `users` (Password Bcrypt cost factor >= 12).
- [x] Seed 1 Profil Staf Default di `staff_profiles` (`user_id = super_admin.id`).
- [x] Seed penugasan cabang default di `user_branches`.

---

### Fase 5.5: Dokumentasi API Interaktif (OpenAPI 3.0.3 & Swagger UI)
- [x] Siapkan loader embedded specs di `docs/docs.go` menggunakan `embed.FS`.
- [x] Buat halaman interaktif Swagger UI di `docs/index.html` dengan dropdown multi-spesifikasi.
- [x] Susun Master OpenAPI 3.0.3 YAML di `docs/swagger.yaml`.
- [x] Susun spesifikasi modular Auth & Multi-Branch di `docs/modules/auth.yaml` (TRD-01).
- [x] Susun spesifikasi modular User & Staff di `docs/modules/user.yaml`.
- [x] Daftarkan route `/docs`, `/docs/swagger.yaml`, `/docs/modules/:name`, `/swagger` di `internal/bootstrap/bootstrap.go`.
- [x] Buat unit testing endpoint dokumentasi di `internal/bootstrap/bootstrap_test.go`.
- [x] Tambahkan skill `.agents/skills/api-docs-generator/SKILL.md` dan anotasi Swaggo pada handler & main.go (merujuk ke pola repo `sehs-be`).

---

### Fase 6: Verifikasi & Proof of Work (PoW)
- [x] Verifikasi build CLI migrasi: `go build ./cmd/migrate`.
- [x] Verifikasi build runner seeder: `go build ./cmd/seed`.
- [x] Jalankan testing: `go test -v ./...` dan pastikan `go build ./cmd/api` sukses tanpa error.
- [x] Buat commit Conventional Commit per domain scope dan siapkan Pull Request.
