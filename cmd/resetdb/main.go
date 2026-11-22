// Resets a local database: drop → migrate → seed.
// Usage: go run ./cmd/resetdb
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

	fmt.Println("→ drop tables...")
	fmt.Println("→ migrate...")
	fmt.Println("→ seed...")

	if err := db.Apply(context.Background(), conn, db.Options{
		DevReset: true,
		Migrate:  true,
		AddSeed:  true,
	}); err != nil {
		log.Fatalf("reset failed: %v", err)
	}

	fmt.Println("✓ done")
}
