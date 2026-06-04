// Reset lokalnej bazy: drop → migrate → seed.
// Użycie: go run ./cmd/resetdb
// Wymaga zmiennej środowiskowej DB_DATABASE_URL.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"uszatakuchnia/db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DB_DATABASE_URL")
	if dsn == "" {
		log.Fatal("DB_DATABASE_URL is not set")
	}

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
