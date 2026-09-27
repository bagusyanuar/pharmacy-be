package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/bagusyanuar/pharmacy-be/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		runMigrateUp()
	case "down":
		runMigrateDown()
	case "step":
		if len(os.Args) < 3 {
			log.Fatalf("usage: go run ./cmd/migrate step <n>")
		}
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid step count: %v", err)
		}
		runMigrateStep(n)
	case "version":
		runMigrateVersion()
	case "force":
		if len(os.Args) < 3 {
			log.Fatalf("usage: go run ./cmd/migrate force <version>")
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid version: %v", err)
		}
		runMigrateForce(v)
	case "create":
		if len(os.Args) < 3 {
			log.Fatalf("usage: go run ./cmd/migrate create <migration_name>")
		}
		name := os.Args[2]
		createMigrationFiles(name)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run ./cmd/migrate <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  up              Run all available up migrations")
	fmt.Println("  down            Roll back the most recent migration (step -1)")
	fmt.Println("  step <n>        Run n up (positive) or down (negative) migrations")
	fmt.Println("  version         Print current migration version")
	fmt.Println("  force <v>       Force set migration version (e.g. after resolving a dirty state)")
	fmt.Println("  create <name>   Create a new pair of .up.sql and .down.sql migration files")
}

func initMigrate() (*migrate.Migrate, *sql.DB) {
	cfg := config.Load()

	dsn := cfg.DSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver instance: %v", err)
	}

	migrationsDir := "file://database/migrations"
	m, err := migrate.NewWithDatabaseInstance(migrationsDir, "postgres", driver)
	if err != nil {
		log.Fatalf("failed to initialize migrate instance: %v", err)
	}

	return m, db
}

func runMigrateUp() {
	m, db := initMigrate()
	defer db.Close()

	log.Println("Running migrations UP...")
	start := time.Now()
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migration changes to apply (database schema is up to date).")
			return
		}
		log.Fatalf("failed to run migrations up: %v", err)
	}

	v, dirty, _ := m.Version()
	log.Printf("Successfully applied migrations UP to version %d (dirty: %t) in %v\n", v, dirty, time.Since(start))
}

func runMigrateDown() {
	m, db := initMigrate()
	defer db.Close()

	log.Println("Running migrations DOWN (1 step)...")
	start := time.Now()
	if err := m.Steps(-1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migration changes to rollback.")
			return
		}
		log.Fatalf("failed to run migration down: %v", err)
	}

	v, dirty, _ := m.Version()
	log.Printf("Successfully rolled back migration. Current version: %d (dirty: %t) in %v\n", v, dirty, time.Since(start))
}

func runMigrateStep(n int) {
	m, db := initMigrate()
	defer db.Close()

	log.Printf("Running migration steps: %d...\n", n)
	start := time.Now()
	if err := m.Steps(n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migration changes.")
			return
		}
		log.Fatalf("failed to run migration step: %v", err)
	}

	v, dirty, _ := m.Version()
	log.Printf("Successfully applied steps. Current version: %d (dirty: %t) in %v\n", v, dirty, time.Since(start))
}

func runMigrateVersion() {
	m, db := initMigrate()
	defer db.Close()

	v, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Println("No migrations have been applied yet (version: nil).")
			return
		}
		log.Fatalf("failed to get migration version: %v", err)
	}

	log.Printf("Current migration version: %d (dirty: %t)\n", v, dirty)
}

func runMigrateForce(version int) {
	m, db := initMigrate()
	defer db.Close()

	log.Printf("Forcing migration version to: %d...\n", version)
	if err := m.Force(version); err != nil {
		log.Fatalf("failed to force migration version: %v", err)
	}
	log.Printf("Successfully forced migration version to %d\n", version)
}

func createMigrationFiles(name string) {
	dir := "database/migrations"
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create migrations directory: %v", err)
	}

	// Count existing migrations to determine next sequential number
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read migrations directory: %v", err)
	}

	highestSeq := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if len(filename) >= 6 {
			if seq, err := strconv.Atoi(filename[:6]); err == nil && seq > highestSeq {
				highestSeq = seq
			}
		}
	}

	nextSeq := highestSeq + 1
	baseFilename := fmt.Sprintf("%06d_%s", nextSeq, name)

	upPath := filepath.Join(dir, baseFilename+".up.sql")
	downPath := filepath.Join(dir, baseFilename+".down.sql")

	if err := os.WriteFile(upPath, []byte("-- Migration UP\n"), 0644); err != nil {
		log.Fatalf("failed to create up migration file: %v", err)
	}
	if err := os.WriteFile(downPath, []byte("-- Migration DOWN\n"), 0644); err != nil {
		log.Fatalf("failed to create down migration file: %v", err)
	}

	fmt.Printf("Created migration files:\n  %s\n  %s\n", upPath, downPath)
}
