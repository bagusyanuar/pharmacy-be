---
name: ticket-clarifier
description: >-
  Facilitates technical clarification, brainstorming, and architectural decision-making for task tickets in bagusyanuar/pharmacy-docs.
  Identifies technical ambiguities, library options, implementation methods, and edge cases, presenting them
  as structured trade-offs with concrete options (A vs B). Trigger before executing any ticket or whenever
  the user wants to brainstorm, clarify, or align on implementation details.
---

# Ticket Clarifier & Brainstorming Partner

This skill governs how the AI agent identifies technical ambiguities, library choices, implementation methods, and facilitates brainstorming for tickets from **`bagusyanuar/pharmacy-docs`**.

---

## 🧭 Core Philosophy: Zero Unilateral Assumptions

1. **User Authority:** The user is the Lead Architect and Decision Maker. The agent NEVER makes unilateral decisions regarding library dependencies, architectural strategies, or database execution methods.
2. **No Blind Execution:** Before writing any code or modifying dependencies, all technical ambiguities and trade-offs must be clarified and approved.
3. **Structured Sparring:** During brainstorming, the agent acts as a sharp architectural advisor—presenting well-reasoned options with pros and cons, but leaving the final choice to the user.

---

## 🎯 When to Activate

This skill MUST be activated in two scenarios:

1. **Automatic Pre-Execution Gate:** Immediately after fetching a ticket and creating the work branch, before any code or migration file is created.
2. **Interactive Brainstorming Mode:** Whenever the user asks to brainstorm, explore approaches, or discuss a ticket (e.g. *"mari kita brainstorm tiket #1"*, *"menurutmu gimana pendekatan terbaik untuk tiket ini?"*).

---

## 🔍 Clarification Categories & Checklist

Before executing a ticket, analyze and evaluate the following 3 core areas:

### 1. Library & External Dependencies
* Inspect existing dependencies in `go.mod`.
* If a new dependency is contemplated (e.g. crypto/hashing, UUID, validation, router/middleware):
  * **Never run `go get` without explicit confirmation.**
  * Clarify whether to use an established standard Go library or a specific third-party package.

### 2. Implementation Methods & Technical Architecture
* **Database & Migration Strategy:**
  * Raw SQL migration files (e.g. `database/migrations/*.sql`) vs GORM AutoMigrate vs CLI migration runners.
* **Transaction & Concurrency Scoping:**
  * Where does database transaction management live (usecase vs infrastructure)?
  * How are race conditions or locks handled (e.g. `SELECT FOR UPDATE`)?
* **Domain & Module Boundaries:**
  * Which DDD module owns each entity or table?
  * How are cross-module dependencies composed (ensuring batch fetching `GetByIDs` instead of GORM Preload)?

### 3. Business Edge Cases & Ambiguities
* Identify gaps in the PRD/TRD/DRA documents.
* Format validations (e.g., NIK 16 digit, phone numbers, SIPA expiration rules).
* Partial indexing rules (`WHERE deleted_at IS NULL`) for soft-deleted entities.

---

## 📝 Standard Presentation Format (Options & Trade-Offs)

When presenting open questions or brainstorming results, ALWAYS use the following structured format:

```markdown
### 🎯 Keputusan #[N]: [Judul Keputusan / Titik Desain]

**Latar Belakang / Konteks:**
[Penjelasan singkat mengapa keputusan ini krusial untuk tiket ini]

* **Opsi A (Rekomendasi): [Nama Pendekatan]**
  * **Kelebihan:** [Poin kelebihan]
  * **Kekurangan / Trade-off:** [Poin konsekuensi/kompleksitas]
  * **Implementasi Teknis:** [Gambaran singkat cara kerjanya]

* **Opsi B: [Nama Pendekatan Alternatif]**
  * **Kelebihan:** [Poin kelebihan]
  * **Kekurangan / Trade-off:** [Poin konsekuensi/kompleksitas]
  * **Implementasi Teknis:** [Gambaran singkat cara kerjanya]

👉 **Rekomendasi Agent:** [Pilihan agent beserta alasan arsitekturalnya]
❓ **Keputusan Kamu:** Apakah kamu setuju dengan Opsi A, atau lebih memilih Opsi B?
```

---

## 💡 Brainstorming Mode Workflow

When the user asks to brainstorm a ticket:

1. **Deconstruct the Ticket:** Break the ticket into architectural sub-problems (Schema, Domain Invariants, Business Logic, API Contract).
2. **Surface Trade-offs:** Present alternative patterns with pros/cons strictly adhering to `.agents/rules/` (Clean Architecture, DDD, and Pharmacy SSOT).
3. **Capture & Summarize Decisions:** Once the user makes decisions, compile an **Implementation Blueprint** summarizing:
   * Approved libraries
   * Approved directory & file structures
   * Approved methods/patterns
4. **Transition to Execution:** Prompt the user if they are ready to proceed with implementation based on the agreed blueprint.
