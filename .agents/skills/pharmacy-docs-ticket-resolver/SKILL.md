---
name: pharmacy-docs-ticket-resolver
description: >-
  Resolves tasks and tickets from the Single Source of Truth repository (bagusyanuar/pharmacy-docs).
  Fetches ticket details via GitHub CLI (gh), extracts SSOT references (PRD, DRA, TRD, DBML),
  decomposes deliverables across Clean Architecture + DDD layers, enforces pharmaceutical domain
  invariants, verifies implementation with tests, and formats Pull Requests linked to the ticket.
  Trigger whenever the user mentions working on, implementing, or resolving a ticket from pharmacy-docs.
---

# Pharmacy Docs Ticket Resolver

This skill guides the AI agent in resolving tasks and implementing features driven by work tickets in the Single Source of Truth (SSOT) repository: **`bagusyanuar/pharmacy-docs`**.

---

## 🧭 Philosophy & Architecture Contract

1. **SSOT Authority:** All business rules (PRD), database schemas (DRA/DBML), and API specifications (TRD) live in `bagusyanuar/pharmacy-docs`. Code in `pharmacy-be` must strictly conform to these specifications.
2. **Ticket-Driven Development:** Every PR in `pharmacy-be` must resolve a specific ticket in `pharmacy-docs` using the syntax `Closes bagusyanuar/pharmacy-docs#<ID>`.
3. **Automated Feedback Loop:** When a PR is merged into `main`, GitHub Actions (`sync-docs-issue.yml`) automatically posts the PR status, merge commit, and documented edge cases back to the referenced ticket in `pharmacy-docs`.

---

## 🛠️ Step-by-Step Resolution Workflow

### Step 1: Fetch and Inspect the Ticket

When given a ticket number (e.g. `Tiket #1` or `tiket 1`):

1. **Query the ticket via `gh CLI`:**
   ```bash
   gh issue view <TICKET_ID> --repo bagusyanuar/pharmacy-docs --json number,title,body,labels
   ```
2. **Identify the standard ticket components:**
   * **📌 Ringkasan Tugas:** High-level summary of the requirement.
   * **📚 Dokumen Acuan (SSOT):** Direct links/paths to PRD, DRA, TRD, or DBML files in `pharmacy-docs`.
   * **🎯 Cakupan Tugas & Checklist Deliverable:** The explicit checklist of items that must be produced.
   * **🛡️ Business Rules yang Wajib Dipenuhi:** Domain constraints (e.g., `BR-PHARM-AUTH-13`, `BR-PHARM-AUTH-02`).
   * **⚠️ Perhatian Khusus & Invarian:** Naming rules, audit trails, precision standards, and UUID requirements.
3. **Create Dedicated Work Branch (via `ticket-branch-creator` skill):**
   * Follow `.agents/skills/ticket-branch-creator/SKILL.md`.
   * Format: `<type>/docs-<id>-<short-slug>` (e.g., `feat/docs-1-auth-schema-migration`, `fix/docs-3-prescription-discount`).
   * Command:
     ```bash
     git checkout main && git pull origin main
     git checkout -b <type>/docs-<id>-<short-slug>
     ```

---

### Step 2: Read Referenced SSOT Documents

Tickets often reference specific architectural docs in `pharmacy-docs`. If detailed specifications (e.g., exact DDL, enum values, DTO payload schemas) are needed:

1. **Read file contents from `pharmacy-docs` using `gh api`:**
   ```bash
   gh api repos/bagusyanuar/pharmacy-docs/contents/<file_path> --jq '.content' | base64 -d
   ```
   *Common paths in `pharmacy-docs`:*
   * Database Architecture / DDL: `technical/01-dra-database-erd-master-auth.md`
   * ERD / DBML Schema: `technical/database/schema.dbml`
   * API TRD: `technical/01-trd-auth-dan-cabang.md`
   * Business PRD: `features/01-prd-auth-user.md`

---

---

### Step 3: Technical Alignment & Open Questions (via `ticket-clarifier` skill)

**MANDATORY GATE:** Before creating files or writing code, execute `.agents/skills/ticket-clarifier/SKILL.md`:
1. Inspect `go.mod` and identify any required libraries (crypto, uuid, validator, etc.).
2. Evaluate technical methods (migration strategies, transaction handling, domain boundaries).
3. Present open questions with structured options (Option A vs Option B) and trade-offs.
4. **STOP and WAIT** for explicit user decision. Do NOT proceed to implementation on unconfirmed assumptions.

---

### Step 4: Decompose Deliverables into Clean Architecture Layers

Map the checklist items to the appropriate layers in `pharmacy-be`:

1. **Database Migrations & Schemas:**
   * Tables with UUID primary keys (`DEFAULT gen_random_uuid()`).
   * Universal Audit Trail (5 columns: `id`, `created_at`, `updated_at`, `created_by`, `deleted_at`).
   * Foreign keys with explicit `ON DELETE RESTRICT` or `ON DELETE SET NULL` constraints.
   * Partial indexes for soft delete: `WHERE deleted_at IS NULL`.
2. **Layer 1: Domain (`internal/modules/<domain>/domain/`):**
   * Pure Go domain entities, value objects, and repository interfaces.
   * Sentinel errors (e.g., `ErrUserNotFound`).
   * **Zero Fiber/GORM imports.**
3. **Layer 2: Usecase (`internal/modules/<domain>/usecase/`):**
   * Business rules enforcement (e.g., password hashing with Bcrypt cost >= 12, supervisor override rules).
   * Cross-module boundaries via foreign repository interfaces (never GORM preloads across modules).
4. **Layer 3: Infrastructure (`internal/modules/<domain>/infrastructure/`):**
   * GORM repository implementations (`*gorm.DB`).
   * Query optimizations, pagination, error translations.
5. **Layer 4: Delivery (`internal/modules/<domain>/delivery/`):**
   * Request/Response DTOs with `validate:"..."` tags.
   * Fiber route handlers using `response.ValidateOrFail` and `response.Success` / `response.Fail`.
6. **Dependency Wiring (`internal/bootstrap/`):**
   * Register new repositories, usecases, and handlers into `SetupApp`.
7. **Seeders / Fixtures (`database/seeders/`):**
   * Pre-seed default master data (e.g., standard roles, default branch, initial admin credentials).

---

### Step 5: Validate Domain & Architectural Invariants

Before finalizing, verify all changes against the core architectural invariants defined in `.agents/rules/pharmacy-docs-ssot.md` and `.agents/rules/architecture.md`:

- [ ] **Multi-Branch Isolation:** Every branch-specific entity has `branch_id UUID NOT NULL REFERENCES branches(id)`.
- [ ] **Database Ubiquitous Language:** All tables/columns in `snake_case` English, except standard Indonesian regulatory terms (`sipa_number`, `strttk_number`, `is_apa`, `sia_number`, `satusehat_ihs_id`, `bpjs_card_number`, `nik`, `tuslah_fee`, `embalase_fee`).
- [ ] **Universal Audit Trail:** 5 standard columns on every business entity table (`id`, `created_at`, `updated_at`, `created_by`, `deleted_at`).
- [ ] **Financial & Metric Precision:**
  * Currency/Monetary: `DECIMAL(12,2)` / int64 (zero floating point loss).
  * Medication dosage & compounding units: `DECIMAL(10,3)`.
- [ ] **Decoupling IAM vs Staff (BR-PHARM-AUTH-13):**
  * `users` = independent IAM login credentials.
  * `staff_profiles` = physical HR employee master, `user_id UUID UNIQUE NULL REFERENCES users(id) ON DELETE SET NULL`.
- [ ] **Zero Hard-Delete (BR-PHARM-AUTH-14):**
  * Soft delete via `deleted_at TIMESTAMPTZ NULL` to preserve historical audit logs and prescription records.
- [ ] **Credential Security:** Bcrypt cost factor >= 12 for passwords.

---

### Step 6: Verification & Proof of Work (PoW)

1. **Tidy dependencies:**
   ```bash
   go mod tidy
   ```
2. **Execute tests:**
   ```bash
   go test -v ./...
   ```
3. **Verify build:**
   ```bash
   go build ./cmd/api
   ```
4. **Collect Proof of Work:**
   * Capture test run summaries and migration/API output to include in the PR description.

---

### Step 7: Commit and PR Preparation (Fully Automated by Agent)

1. **Commit Changes:**
   * Use the `conventional-commit` skill (`<type>(<domain>): <summary>`).
   * Confirm with the user before committing.
2. **Auto-Generate Pull Request Description:**
   * The AI Agent automatically prepares the complete markdown body adhering to [`.github/pull_request_template.md`](file:///Users/dystopia/projects/pharmacy/pharmacy-be/.github/pull_request_template.md):
     * **Link:** `Closes bagusyanuar/pharmacy-docs#<TICKET_ID>`.
     * **Description:** Concise summary of all technical deliverables completed.
     * **Type of Change:** Check the matching box (`[x]`).
     * **Compliance Checklist:** Check off all verified architecture rules (`[x]`).
     * **Proof of Work:** Automatically embed the captured `go test` and migration terminal logs into the code blocks.
     * **💡 Pekerjaan Tambahan / Edge Cases:** Automatically formulate concise bullet points of any additional findings, edge cases, or extra fixes made outside the initial scope.
3. **Submit / Open PR via `gh pr create`:**
   * Propose the exact command and show the complete PR body draft to the user for review:
     ```bash
     gh pr create --base main --head <branch-name> --title "<type>(<domain>): <title>" --body-file pr_body.md
     ```
   * Once merged, `sync-docs-issue.yml` will automatically sync the status and edge cases back to the docs ticket!
