---
name: api-docs-generator
description: >-
  Generate and maintain OpenAPI / Swagger documentation for Fiber v3 REST APIs.
  Use when adding swagger annotations to handlers, writing or updating OpenAPI 3.0.3 YAML specs,
  configuring Swagger UI, or adding new module API contracts.
---

# API Documentation Generator (Swagger / OpenAPI)

This skill guides generating, writing, and maintaining OpenAPI 3.0.3 specifications and interactive Swagger UI documentation for Fiber v3 REST APIs in this repository, following the pattern established in `smart-environment-health-system/sehs-be`.

---

## 1. Documentation Architecture

Documentation is embedded directly into the Go binary using `embed.FS` and served via Fiber v3 routes:

```text
docs/
├── docs.go                  # embed.FS embedding index.html, swagger.yaml, modules/*.yaml
├── index.html               # Custom-styled Swagger UI with multi-spec dropdown selector
├── swagger.yaml             # Master / Combined API specification (all modules)
├── modules/                 # Modular OpenAPI 3.0.3 YAML specs per domain context
│   ├── auth.yaml            # Auth, POS Fast-Login, Session, Supervisor Override (TRD-01)
│   └── user.yaml            # IAM Accounts & Profiles
└── README.md                # Guide on viewing and importing specs (Postman/Insomnia/Editor)
```

### Endpoints Exposed by Fiber (`internal/bootstrap/bootstrap.go`):
- `GET /docs`: Interactive Swagger UI with brand styling and spec switcher dropdown.
- `GET /swagger`, `GET /swagger/*`, `GET /api-docs`: 302 Redirect to `/docs`.
- `GET /docs/swagger.yaml`: Master OpenAPI YAML specification.
- `GET /docs/modules/:name`: Granular OpenAPI YAML per module (e.g. `/docs/modules/auth.yaml`).

---

## 2. General API Annotations (`cmd/api/main.go`)

When using code annotations with swaggo:

```go
// @title           Pharmacy Information System (POS & ERP) API
// @version         1.0
// @description     High-performance multi-branch pharmacy management REST API built with Go, Fiber v3, GORM, and PostgreSQL.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Pharmacy Engineering
// @contact.url     https://github.com/bagusyanuar/pharmacy-be

// @license.name    Proprietary
// @license.url     https://github.com/bagusyanuar/pharmacy-be

// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and your JWT token.
```

---

## 3. Handler Method Annotations

Every Fiber handler must be decorated with standard swaggo annotations conforming to `pkg/response`:

### 3.1 Standard GET (Paginated List) Example
```go
// List godoc
// @Summary      List users
// @Description  Get a paginated list of registered users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page      query     int     false  "Page number (default 1)"   default(1)
// @Param        per_page  query     int     false  "Items per page (default 20, max 100)" default(20)
// @Success      200       {object}  response.Response{data=[]usecase.UserResponseDTO,meta=response.PaginationMeta}
// @Failure      401       {object}  response.Response
// @Failure      500       {object}  response.Response
// @Router       /users [get]
func (h *UserHandler) List(c fiber.Ctx) error {
    // ...
}
```

### 3.2 Standard POST (Creation / Mutation) Example
```go
// Register godoc
// @Summary      Register new user
// @Description  Register a new user account with role binding
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      delivery.RegisterRequest  true  "Register Request Payload"
// @Success      201      {object}  response.Response{data=usecase.UserResponseDTO}
// @Failure      400      {object}  response.Response  "Invalid request body format"
// @Failure      422      {object}  response.Response  "Validation error with field details"
// @Failure      409      {object}  response.Response  "Email already registered"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
    // ...
}
```

---

## 4. Modular YAML Conventions (`docs/modules/*.yaml`)

For each bounded context / ticket, maintain a dedicated YAML specification:
1. **SSOT Alignment**: Reference the exact ticket and specification (e.g. `TRD-PHARM-01`).
2. **Standard Response Envelopes**:
   - `SuccessEnvelope`: `{ success: true, message: string, data: object }`
   - `ErrorEnvelope`: `{ success: false, error: { code: string, message: string, details: array|null } }`
   - `PaginatedResponse`: `{ success: true, message: string, data: array, meta: { current_page, per_page, total_items, total_pages } }`
3. **Register in `docs/index.html`**:
   Add the new module to the Swagger UI `urls` dropdown list:
   ```javascript
   urls: [
     { url: "/docs/swagger.yaml", name: "🌟 Master Spec (All Features)" },
     { url: "/docs/modules/auth.yaml", name: "🔐 Auth & Multi-Branch Module (TRD-01)" },
     { url: "/docs/modules/branch.yaml", name: "🏢 Branch Management Module" },
     ...
   ]
   ```

---

## 5. Documentation Checklist

- [ ] Does `docs/docs.go` embed all YAML files (`//go:embed index.html swagger.yaml modules/*.yaml`)?
- [ ] Are `/docs`, `/swagger`, `/api-docs`, `/docs/swagger.yaml`, and `/docs/modules/:name` registered in `internal/bootstrap/bootstrap.go`?
- [ ] Are test assertions in `internal/bootstrap/bootstrap_test.go` passing for all doc routes?
- [ ] Does every endpoint document `200/201`, `400/422`, `401/403`, and `404/500` scenarios?
- [ ] Are protected endpoints marked with `security: - BearerAuth: []`?
