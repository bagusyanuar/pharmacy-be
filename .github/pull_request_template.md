## 🔗 Tiket Terkait (Docs SSOT)
<!-- Cantumkan nomor tiket di repositori pharmacy-docs. Format: Closes bagusyanuar/pharmacy-docs#<ID> -->
Closes bagusyanuar/pharmacy-docs#

---

## 📌 Deskripsi Perubahan
<!-- Rangkuman singkat mengenai pekerjaan yang diselesaikan dalam PR ini -->

---

## 🏷️ Tipe Perubahan
- [ ] 🚀 Fitur Baru (Endpoint API / Service)
- [ ] 🗄️ Migrasi Database / DDL Skema Relasional
- [ ] 🐛 Bugfix / Perbaikan Logika Bisnis
- [ ] ⚡ Optimasi Performa / Indeks
- [ ] 🔧 Refactoring / Konfigurasi CI/CD

---

## ✅ Checklist Kepatuhan Arsitektur Farmasi (Wajib Terpenuhi)
- [ ] **Multi-Branch Isolation:** Validasi `branch_id` konsisten atau middleware `X-Branch-Id` diimplementasikan dengan benar.
- [ ] **Standar Penamaan Database:** Kolom/tabel menggunakan bahasa Inggris (`snake_case`). Istilah regulasi farmasi Indonesia baku dipertahankan (`sipa_number`, `strttk_number`, `is_apa`, `sia_number`, `satusehat_ihs_id`, `bpjs_card_number`, `nik`).
- [ ] **Universal Audit Trail:** 5 kolom audit (`id`, `created_at`, `updated_at`, `created_by`, `deleted_at`) ada pada tabel entitas.
- [ ] **Presisi Finansial & Satuan:** Nilai moneter/harga menggunakan `DECIMAL(12,2)` / int64 (tanpa floating point loss), dosis/satuan racikan menggunakan `DECIMAL(10,3)`.
- [ ] **Decoupling IAM vs Staf (BR-PHARM-AUTH-13):** Tabel `users` adalah principal IAM mandiri; tabel `staff_profiles` memegang `user_id UUID UNIQUE NULL REFERENCES users(id) ON DELETE SET NULL`.
- [ ] **Envelope Response API Standar:** Format output JSON seragam `{ "success": true, "data": ..., "meta": ..., "error": ... }`.

---

## 📸 Bukti Pengerjaan (Proof of Work) — WAJIB DILAMPIRKAN
<!-- Lampirkan bukti log terminal migrasi, output unit test, atau screenshot/log cURL/Postman response -->
- [ ] **Hasil Eksekusi Unit/Integration Test (`make test`):**
  ```text
  (Tempelkan ringkasan log test di sini)
  ```
- [ ] **Hasil Pengujian API / Database Migration (Jika Relevan):**
  ```json
  (Tempelkan payload respons pengujian di sini)
  ```

---

## 💡 Pekerjaan Tambahan / Edge Cases (Otomatis Disinkronkan ke Issue Docs)
<!-- Tuliskan jika ada temuan kasus tepi, bugfix tambahan, atau penanganan khusus di luar tiket awal. Catatan ini akan otomatis diposting oleh bot ke Issue Docs terkait saat PR di-merge. -->
- [ ] ...
