---
name: ticket-todo-generator
description: >-
  Generates a structured technical Todo List markdown document in docs/todo/ based on a task ticket
  from bagusyanuar/pharmacy-docs. Deconstructs ticket requirements across Clean Architecture & DDD layers,
  database migrations, business rules, and verification checklists. Trigger whenever the user asks to create
  a todo list, plan execution, or scaffold tasks for a ticket.
---

# Ticket Todo List Generator

This skill standardizes the creation of actionable, phased technical **Todo Lists** stored in **`docs/todo/`** for tickets originating from **`bagusyanuar/pharmacy-docs`**.

---

## 🏷️ File Naming Convention

All generated todo list documents must be stored in:

$$\mathbf{docs/todo/ticket-\langle id \rangle -\langle short-slug \rangle .md}$$

* Example: `docs/todo/ticket-01-auth-schema-migration-and-seeding.md`
* Format: All lowercase, numbers zero-padded if single digit (`ticket-01`), hyphen-separated `kebab-case`.

---

## 🛠️ Step-by-Step Workflow

When requested to generate a todo list for a ticket (e.g. *"buatkan todo list untuk tiket #2"*):

### Step 1: Fetch Ticket Metadata
Query ticket information from `bagusyanuar/pharmacy-docs` using GitHub CLI:
```bash
gh issue view <TICKET_ID> --repo bagusyanuar/pharmacy-docs --json number,title,body,labels
```

### Step 2: Extract Key Metadata
* **Ticket ID & Title:** e.g. `[BE] Auth & Sesi (Part 2/5): Endpoint Web Login, Fast PIN POS...`
* **Target Branch:** Derived via `ticket-branch-creator` (e.g. `feat/docs-2-auth-web-login-fast-pin`).
* **Acuan SSOT:** Relevant PRD, TRD, DRA links cited in the ticket.

### Step 3: Deconstruct into Phased Technical Checklist
Structure the todo document using standard phased sections:

```markdown
# 📋 Todo List: Tiket #[ID] — [Judul Tiket]

* **Tiket Acuan:** [bagusyanuar/pharmacy-docs#[ID]](https://github.com/bagusyanuar/pharmacy-docs/issues/[ID])
* **Modul:** [Nama Modul]
* **Target Branch:** `[type]/docs-[id]-[slug]`
* **Status:** 🟡 Dalam Perencanaan / Siap Eksekusi

---

## 🎯 Deliverable Checklist

### Fase 1: Persiapan & Skema Database (Jika Relevan)
- [ ] File migrasi SQL / DDL ...
- [ ] Penegakan constraint & index ...

### Fase 2: Layer Domain (`internal/modules/<domain>/domain/`)
- [ ] Pure entities & value objects ...
- [ ] Repository interfaces ...
- [ ] Sentinel errors di `errors.go` ...

### Fase 3: Layer Usecase (`internal/modules/<domain>/usecase/`)
- [ ] Business logic & orchestration ...
- [ ] Cross-module repository dependency injection ...

### Fase 4: Layer Infrastructure (`internal/modules/<domain>/infrastructure/`)
- [ ] GORM repository implementation ...
- [ ] Query optimization & error translation ...

### Fase 5: Layer Delivery (`internal/modules/<domain>/delivery/`)
- [ ] Request & Response DTOs with `validate:"..."` tags ...
- [ ] Fiber v3 Route Handlers & error mapping ...

### Fase 6: DI Wiring di Bootstrap (`internal/bootstrap/bootstrap.go`)
- [ ] Registrasi repository, usecase, dan handler baru ...

### Fase 7: Pengujian & Proof of Work (PoW)
- [ ] Unit & integration tests (`make test`) ...
- [ ] Verifikasi build (`go build ./cmd/api`) ...
- [ ] Penyusunan Pull Request draft ...
```

### Step 4: Write to `docs/todo/`
Save the generated content to `docs/todo/ticket-<id>-<short-slug>.md`.

### Step 5: Active Progress Tracking
During implementation of the ticket, the AI agent updates the `- [ ]` checkboxes in this file to `- [x]` as deliverables are completed, maintaining a transparent, persistent progress log.
