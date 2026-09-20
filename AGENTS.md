# Golang Clean Architecture & DDD Backend Template

This project follows Clean Architecture and Domain-Driven Design (DDD) principles using Golang, Fiber v3, GORM, and PostgreSQL.

## Architecture Overview

```text
├── cmd
│   └── api
│       └── main.go                 # Entry point: process lifecycle only (signals, graceful shutdown)
├── internal
│   ├── bootstrap                   # Central DI wiring: Repository -> Usecase -> Handler, per module
│   ├── config                      # Application configuration (.env, database setup)
│   ├── middleware                  # Fiber middlewares (Auth, Logger, CORS)
│   └── modules                     # Domain/Context separation (DDD)
│       └── [module_name]           # Feature / Domain Module
│           ├── domain              # Layer 1: Core Business (Entities, Repo Interfaces, Errors)
│           ├── usecase             # Layer 2: Application Business Rules & DTOs
│           ├── infrastructure      # Layer 3: GORM implementation, External APIs
│           └── delivery            # Layer 4: Fiber HTTP Handlers & Routes
├── pkg                             # Global helpers/utilities (logger, response, validator, jwt)
└── scripts                         # Scaffolding & project initialization utilities
```

## Critical Agent Instructions

Whenever you are writing, refactoring, or generating code in this repository:
1. **Follow the Rules in `.agents/rules/`**: Read and strictly adhere to all guidelines.
2. **Domain Layer Independence**: The `domain/` layer must NEVER import `github.com/gofiber/fiber/v3` or `gorm.io/gorm`. It must only contain pure Go domain models, interfaces, and sentinel errors.
3. **Cross-Module Boundaries**: NEVER use cross-module GORM relations or `Preload()` across modules. Usecases must depend on foreign repository interfaces and use minimal summary DTOs with batch fetching (`GetByIDs`).
4. **Use Skill for New Modules**: When asked to add a new entity, module, or feature, execute according to `.agents/skills/go-module-generator/SKILL.md`.
5. **Use Skill for Commits**: When asked to commit changes or create commit messages, execute according to `.agents/skills/conventional-commit/SKILL.md`.

## Rules Index

- [architecture.md](.agents/rules/architecture.md): Complete guidelines for Directory Structure, Layers, Dependency Inversion, and Cross-Module Composition.
- [naming-convention.md](.agents/rules/naming-convention.md): Rules for naming variables, structs, files, and interfaces across all layers.
- [error-handling.md](.agents/rules/error-handling.md): Guidelines for custom sentinel errors, wrapping, and HTTP status code mapping.
- [data-modeling.md](.agents/rules/data-modeling.md): Entity/schema design conventions (UUID primary keys, soft delete, timestamps).
- [query-optimization.md](.agents/rules/query-optimization.md): GORM query best practices — avoiding N+1, pagination, indexing, and batch fetching.
- [logging.md](.agents/rules/logging.md): Structured logging with Zap — dependency injection, log levels, contextual fields.
- [response.md](.agents/rules/response.md): Standardized JSON response envelope (`pkg/response`) — success, failure, and pagination meta.
- [validation.md](.agents/rules/validation.md): Request DTO validation conventions using `pkg/validator` + `response.ValidateOrFail`.
- [authentication.md](.agents/rules/authentication.md): JWT access & refresh token pattern, HttpOnly cookie handling, and middleware protection.
