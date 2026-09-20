---
name: ticket-branch-creator
description: >-
  Creates and checks out Git feature/fix branches mapped directly to task tickets in bagusyanuar/pharmacy-docs.
  Ensures the branch is cut from an up-to-date main branch using the standard naming format <type>/docs-<id>-<slug>.
  Trigger whenever the user asks to create a branch, start a new ticket, or checkout a branch for a ticket.
---

# Ticket Branch Creator

This skill standardizes the creation and checkout of Git branches for tasks originating from work tickets in **`bagusyanuar/pharmacy-docs`**.

---

## 🏷️ Branch Naming Standard

All branches must adhere strictly to the following format:

$$\mathbf{\langle type \rangle /docs-\langle id \rangle -\langle short-slug \rangle}$$

### 1. Type Prefix Mapping
Select the prefix matching the nature of the ticket:
* `feat/` — New feature, relational schema migration, DDL changes, or new API endpoints.
* `fix/` — Bug fixes or business logic corrections.
* `perf/` — Query tuning, indexing, caching, or performance optimizations.
* `refactor/` — Code restructuring or cleanup without behavior changes.
* `chore/` — CI/CD, dependency updates, tooling, or project configuration.

### 2. Ticket Identifier (`docs-<id>`)
* Always use `docs-<id>` (e.g. `docs-1`, `docs-14`).
* **FORBIDDEN:** Do NOT use `#` in branch names (e.g. `feat/#1` or `feat/ticket#1`), as `#` causes syntax errors and accidental comments in bash and zsh shells.

### 3. Slug (`<short-slug>`)
* 2–4 descriptive words in English.
* Format: `kebab-case` (all lowercase, hyphen-separated).
* Derived directly from the ticket title.

---

## 🛠️ Execution Workflow

When requested to create a branch for a ticket (e.g. *"buatkan branch untuk tiket #1"*):

### Step 1: Inspect Ticket Metadata
Fetch the ticket title and labels using GitHub CLI:
```bash
gh issue view <TICKET_ID> --repo bagusyanuar/pharmacy-docs --json number,title,labels
```

*Example:*
- Tiket #1: `[BE] Auth & Sesi (Part 1/5): Migrasi Skema Database Relasional, Indeks & Seeding`
- Type: `feat`
- Slug: `auth-schema-migration`
- Target Branch Name: `feat/docs-1-auth-schema-migration`

### Step 2: Safety Check on Working Tree
Ensure the current working tree is clean before switching branches:
```bash
git status
```
If uncommitted files exist:
- Alert the user and resolve (stash or commit) before switching.

### Step 3: Synchronize with Remote `main`
Always cut the new branch from the latest commit on `origin/main`:
```bash
git checkout main
git pull origin main
```

### Step 4: Create and Checkout Branch
Create and switch to the new feature branch:
```bash
git checkout -b <branch_name>
```

*Example:*
```bash
git checkout -b feat/docs-1-auth-schema-migration
```

### Step 5: Verification
Confirm the branch was successfully created and checked out:
```bash
git branch --show-current
```
Report the active branch clearly to the user.
