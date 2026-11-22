// Adds missing columns and tables without deleting data.
// Usage: go run ./cmd/migrate
// Requires the DB_DATABASE_URL environment variable.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"uszatakuchnia/db"
)

func main() {
	dsn := os.Getenv("DB_DATABASE_URL")
	if dsn == "" {
		log.Fatal("DB_DATABASE_URL is not set")
	}

	conn, err := db.Open(dsn, nil)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}

	fmt.Println("→ migrate (only adds missing columns/tables, no data loss)...")

	if err := db.Apply(context.Background(), conn, db.Options{
		DevReset: false,
		Migrate:  true,
		AddSeed:  false,
	}); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	fmt.Println("✓ done")
}
