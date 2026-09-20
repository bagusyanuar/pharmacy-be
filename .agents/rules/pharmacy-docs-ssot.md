# Aturan SSOT Dokumentasi Domain & Spesifikasi Teknis Farmasi

Repositori ini (**`pharmacy-be`**) merupakan repositori implementasi teknis Backend untuk Sistem Informasi Apotek (POS & ERP).

Seluruh persyaratan bisnis (PRD), rancangan arsitektur data (DRA), dan spesifikasi kontrak API (TRD) bersumber secara mutlak dari repositori Single Source of Truth (SSOT):
👉 **`https://github.com/bagusyanuar/pharmacy-docs`**

---

## 📚 Dokumen Acuan Wajib Sebelum Implementasi

1. **Master Data Architecture & Database Relasional (DRA):**
   * File acuan: `technical/01-dra-database-erd-master-auth.md` di `pharmacy-docs`.
   * Memuat skema PostgreSQL 15+ DDL lengkap, tipe data UUID, presisi `DECIMAL`, indeks parsial (`WHERE deleted_at IS NULL`), dan data seeding awal.
2. **Spesifikasi Teknis API & Logika Autentikasi (TRD):**
   * File acuan: `technical/01-trd-auth-dan-cabang.md` di `pharmacy-docs`.
   * Memuat kontrak endpoint RESTful, DTO JSON payload, siklus hidup token (RTR), dan logika otorisasi *Supervisor Override*.
3. **Fondasi Multi-Cabang & Aturan Finansial:**
   * File acuan: `technical/00-architecture-and-multibranch-guidelines.md` di `pharmacy-docs`.
   * Memuat aturan isolasi data multi-cabang, strategi indeks FEFO, dan audit trail universal 5 kolom.

---

## 🔒 Invarian Arsitektur yang Tidak Boleh Dilanggar

1. **Konvensi Penamaan Database:**
   * Seluruh nama tabel dan kolom **WAJIB** menggunakan bahasa Inggris dengan format `snake_case` (misal: `branch_id`, `created_at`, `unit_cost_price`).
   * Istilah legalitas, perizinan farmasi resmi, dan regulasi pemerintah Indonesia **TETAP DIPERTAHANKAN** dalam format aslinya:
     * `sipa_number` (Surat Izin Praktik Apoteker)
     * `strttk_number` (Surat Tanda Registrasi Tenaga Teknis Kefarmasian)
     * `sia_number` (Surat Izin Apotek)
     * `is_apa` (Apoteker Pengelola Apotek)
     * `bpjs_card_number`
     * `sipnap_reported_at` (Sistem Pelaporan Narkotika & Psikotropika)
     * `satusehat_ihs_id` (Kemenkes SatuSehat)
     * `nik` (Nomor Induk Kependudukan 16 digit)
2. **Pemisahan Akun IAM vs Profil Staf (BR-PHARM-AUTH-13):**
   * Entitas `users` adalah principal kredensial login IAM mandiri.
   * Entitas `staff_profiles` memegang foreign key opsional `user_id UUID UNIQUE NULL REFERENCES users(id) ON DELETE SET NULL`.
   * Staf non-sistem (kurir, helper) dapat terdaftar di `staff_profiles` dengan `user_id = NULL`.
3. **Isolasi Multi-Cabang Wajib:**
   * Setiap tabel transaksi, inventaris stok, dan mutasi obat **WAJIB** memiliki foreign key `branch_id UUID NOT NULL REFERENCES branches(id)`.
   * Middleware Backend wajib memvalidasi kecocokan `X-Branch-Id` pada setiap request API yang beroperasi di lingkup cabang.
4. **Presisi Finansial & Satuan Farmasi:**
   * Nilai moneter/keuangan menggunakan `DECIMAL(12,2)` / int64 (sen murni tanpa floating point loss).
   * Dosis obat racikan dan konversi fraksi satuan menggunakan presisi 3 desimal `DECIMAL(10,3)`.
5. **Universal Audit Trail 5 Kolom:**
   * Setiap tabel entitas data bisnis wajib menyertakan 5 kolom audit: `id`, `created_at`, `updated_at`, `created_by`, `deleted_at`.

---

## 🔗 Kaitan Pull Request dengan Tiket Docs

Setiap Pull Request di repo `pharmacy-be` **WAJIB** mencantumkan nomor tiket di repositori `pharmacy-docs` pada deskripsinya:
```text
Closes bagusyanuar/pharmacy-docs#<NOMOR_TIKET>
```
Saat PR di-merge, bot CI/CD akan secara otomatis melaporkan status penyelesaian dan mencatat temuan pekerjaan tambahan (*edge cases*) ke tiket dokumentasi terkait.
