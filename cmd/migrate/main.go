// Dodaje brakujące kolumny i tabele bez usuwania danych.
// Użycie: go run ./cmd/migrate
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
