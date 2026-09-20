---
name: backend-code-reviewer
description: >-
  Audits and reviews backend Go code for database query optimization (anti-N+1, indexing, pagination),
  code robustness (nil-pointer prevention, ACID transactions, error wrapping), and Clean Architecture/DDD
  domain invariants. Trigger whenever the user asks to review code, audit queries, check for N+1 issues,
  or verify production-readiness before submitting a PR.
---

# Backend Code Reviewer & Quality Auditor

This skill establishes the standard procedure for conducting comprehensive architectural and performance code reviews in **`pharmacy-be`**, both locally via CLI and automatically via GitHub Actions (Gemini Pro).

---

## 🧭 The 3 Core Review Pillars

Every piece of code submitted or modified must be evaluated against the following 3 pillars:

---

### Pillar 1: Query & Database Performance (Anti N+1 & Indexing)
* **[CRITICAL] Zero N+1 Queries:**
  * **Violation:** Executing database queries inside loops (e.g. `for _, item := range items { db.Where(...).First(&x) }`).
  * **Remedy:** Collect foreign keys in a slice, execute a single batch query `repo.GetByIDs(ctx, ids)`, and construct an in-memory `map[string]*Entity` for $O(1)$ lookup.
* **[CRITICAL] Mandatory Pagination:**
  * All list endpoints **MUST** enforce pagination using `response.ParsePagination(c)` (default 20, max 100).
  * Unbounded `SELECT * FROM table` queries are strictly prohibited.
* **[HIGH] Explicit Column Selection:**
  * Avoid `SELECT *`. Explicitly select required columns with `.Select("id", "name", ...)` to reduce heap allocations and database I/O.
* **[HIGH] Context Cancellation Propagation:**
  * All GORM calls must chain `.WithContext(ctx)` so that client timeouts or cancellations immediately abort the database query.
* **[HIGH] Relational Index Coverage:**
  * Every column used in `WHERE`, `JOIN`, or `ORDER BY` must be backed by a PostgreSQL index (especially composite indexes or partial indexes `WHERE deleted_at IS NULL`).

---

### Pillar 2: Code Robustness & Resilience (Zero Panics & Safe Transactions)
* **[CRITICAL] Nil-Pointer Safety:**
  * Always verify `if entity == nil` before accessing fields or methods on pointer variables.
  * Usecase and repository methods returning pointers must return explicit errors (e.g. `(nil, domain.ErrNotFound)`).
* **[CRITICAL] Safe Database Transactions:**
  * Any operation modifying multiple tables must be wrapped inside a database transaction (`db.Transaction(func(tx *gorm.DB) error { ... })`).
  * Transactions must automatically roll back upon encountering an error or panic.
* **[HIGH] Error Wrapping & Sentinel Errors:**
  * Core domain business errors must be declared in `internal/modules/<domain>/domain/errors.go` (e.g. `var ErrBranchNotFound = errors.New(...)`).
  * Technical errors must be wrapped using `fmt.Errorf("context message: %w", err)`.
  * Handlers must inspect errors using `errors.Is(err, domain.ErrX)`.
* **[HIGH] Input Validation:**
  * Handlers must validate all request DTOs using `response.ValidateOrFail(c, req)` before invoking usecases.

---

### Pillar 3: Clean Architecture & Pharmaceutical Invariants
* **[CRITICAL] Domain Independence:**
  * The `domain/` layer must NEVER import `github.com/gofiber/fiber/v3` or `gorm.io/gorm`.
  * Pure Go structs, interfaces, and standard library types only.
* **[CRITICAL] Multi-Branch Data Isolation:**
  * All operational entities and branch-scoped queries must mandate `branch_id UUID NOT NULL REFERENCES branches(id)`.
  * Handlers must validate the `X-Branch-Id` header.
* **[CRITICAL] Financial & Medication Dosage Precision:**
  * Currency and monetary values: `DECIMAL(12,2)` or integer cents (zero floating point).
  * Medication compounding dosages: `DECIMAL(10,3)`.
* **[HIGH] Universal Audit Trail (5 Columns):**
  * Every entity table must contain: `id`, `created_at`, `updated_at`, `created_by`, `deleted_at`.
* **[HIGH] Decoupling IAM vs Staff (BR-PHARM-AUTH-13):**
  * `users` = independent login credentials.
  * `staff_profiles` = physical HR employee, `user_id UUID UNIQUE NULL REFERENCES users(id) ON DELETE SET NULL`.
* **[HIGH] Uniform JSON Response Envelope:**
  * Endpoints must return responses via `response.Success` / `response.Fail` (`{ success, data, meta, error }`).

---

## 🛠️ Review Execution Workflow

When requested to review code (e.g. *"tolong review code di branch ini"*):

### Step 1: Inspect Changes
Fetch the git diff against `main`:
```bash
git diff origin/main...HEAD
```
Or inspect the staged/modified files:
```bash
git diff --cached
```

### Step 2: Audit Against the 3 Pillars
Examine all modified Go files and SQL migrations against every checklist item.

### Step 3: Produce Structured Scorecard
Format the review output as follows:

```markdown
## 🛡️ Laporan Code Review & Quality Audit

### 📊 Scorecard Ringkasan
| Pilar Audit | Status | Ringkasan Temuan |
| :--- | :---: | :--- |
| **Pilar 1: Query & Performa Database (Anti N+1)** | 🟢 / 🟡 / 🔴 | [Keterangan singkat] |
| **Pilar 2: Ketahanan Kode (Robustness & Nil-Safety)** | 🟢 / 🟡 / 🔴 | [Keterangan singkat] |
| **Pilar 3: Clean Architecture & Invarian Farmasi** | 🟢 / 🟡 / 🔴 | [Keterangan singkat] |

### 🔍 Temuan Mendalam
#### 🔴 Critical Issues (Wajib Diperbaiki Sebelum Merge)
1. **[File:Line]** Deskripsi masalah dan dampaknya.

#### 🟡 Warnings & Rekomendasi Optimasi
1. **[File:Line]** Deskripsi saran perbaikan.

#### 🟢 Positive Highlights
* Aspek kode yang sudah mematuhi best practice.

### 💡 Rekomendasi Perbaikan Kode (Code Diff)
```diff
- // Kode sebelum
+ // Kode sesudah perbaikan
```
```

### Step 4: Offer Immediate Remediation
Ask the user if they would like the agent to automatically apply the recommended fixes.
