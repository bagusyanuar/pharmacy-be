package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/bagusyanuar/pharmacy-be/database/seeders"
	"github.com/bagusyanuar/pharmacy-be/internal/config"
)

func main() {
	cfg := config.Load()

	dsn := cfg.DSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Seeding initial database data...")
	start := time.Now()

	if err := seeders.SeedAuthAndStaff(ctx, db); err != nil {
		log.Fatalf("failed to seed auth and staff: %v", err)
	}

	log.Printf("Successfully completed database seeding in %v\n", time.Since(start))
}
