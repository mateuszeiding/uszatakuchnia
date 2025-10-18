package db

import (
	"fmt"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once     sync.Once
	inst     *gorm.DB
	DevReset = false
)

func DB() *gorm.DB {
	once.Do(func() {
		if err := initDB(); err != nil {
			panic(fmt.Errorf("db init failed: %w", err))
		}
	})

	return inst
}

func initDB() error {
	dsn := os.Getenv("DB_DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DB_DATABASE_URL is empty")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	if DevReset {
		if err := devReset(db); err != nil {
			return fmt.Errorf("dev reset: %w", err)
		}
	}

	if err := runMigrations(db); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	if err := runSeed(db); err != nil {
		return fmt.Errorf("seed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("db.DB(): %w", err)
	}
	sqlDB.SetMaxOpenConns(15)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	inst = db

	return nil
}

func Close() error {
	if inst == nil {
		return nil
	}

	sql, _ := inst.DB()

	return sql.Close()

}
