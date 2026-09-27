# 📚 Dokumentasi API Interaktif (Swagger UI & OpenAPI 3.0.3)

Repositori ini menyediakan dokumentasi API interaktif berbasis **Swagger UI** dan spesifikasi **OpenAPI 3.0.3** yang disematkan (*embedded*) langsung ke dalam binary aplikasi melalui `embed.FS`.

---

## 🚀 Mengakses Dokumentasi

Jalankan server aplikasi:
```bash
make run
# atau: go run ./cmd/api
```

Buka URL berikut di browser:
* 🌐 **Swagger UI Interaktif:** `http://localhost:3000/docs`
* 🔀 **Alias Redirects:** `http://localhost:3000/swagger` atau `http://localhost:3000/api-docs`
* 📄 **Master OpenAPI YAML Spec:** `http://localhost:3000/docs/swagger.yaml`
* 📦 **Auth Module Spec (TRD-01):** `http://localhost:3000/docs/modules/auth.yaml`
* 👥 **User Module Spec:** `http://localhost:3000/docs/modules/user.yaml`

---

## 📂 Struktur File Dokumentasi

```text
docs/
├── docs.go             # Embed binary loader (embed.FS)
├── index.html          # Custom-branded Swagger UI Bundle (CDN)
├── swagger.yaml        # Master OpenAPI 3.0.3 Specification
├── modules/            # Spesifikasi modular per Bounded Context
│   ├── auth.yaml       # Spesifikasi Auth & Multi-Branch (TRD-01)
│   └── user.yaml       # Spesifikasi User & Staff Management
└── README.md           # Panduan dokumentasi ini
```

---

## 🔗 Referensi SSOT
Seluruh kontrak endpoint, skema DTO, dan alur otentikasi merujuk secara mutlak pada repositori Single Source of Truth (SSOT):
👉 [**bagusyanuar/pharmacy-docs**](https://github.com/bagusyanuar/pharmacy-docs)
