package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once sync.Once
	inst *gorm.DB

	DevReset = false
	Migrate  = false
	AddSeed  = false
)

type Options struct {
	DevReset bool
	Migrate  bool
	AddSeed  bool
}

type step struct {
	name    string
	enabled bool
	run     func(tx *gorm.DB) error
}

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

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("db.DB(): %w", err)
	}
	sqlDB.SetMaxOpenConns(15)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	inst = db

	if err := Apply(context.Background(), inst, Options{
		DevReset: false,
		Migrate:  false,
		AddSeed:  false,
	}); err != nil {
		return fmt.Errorf("apply: %w", err)
	}

	return nil
}

func Apply(ctx context.Context, db *gorm.DB, opt Options) error {
	const lockKey int64 = 213742069

	steps := []step{
		{name: "dev reset", enabled: opt.DevReset, run: devReset},
		{name: "migrations", enabled: opt.Migrate, run: runMigrations},
		{name: "seed", enabled: opt.AddSeed, run: runSeed},
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockTx(tx, lockKey); err != nil {
			return err
		}
		return runSteps(tx, steps)
	})
}

func lockTx(tx *gorm.DB, key int64) error {
	if err := tx.Exec(`SELECT pg_advisory_xact_lock(?)`, key).Error; err != nil {
		return fmt.Errorf("advisory lock: %w", err)
	}
	return nil
}

func runSteps(tx *gorm.DB, steps []step) error {
	for _, s := range steps {
		if !s.enabled {
			continue
		}
		if err := s.run(tx); err != nil {
			return fmt.Errorf("%s: %w", s.name, err)
		}
	}
	return nil
}

func MigrateAndSeed(db *gorm.DB) error {
	return Apply(context.Background(), db, Options{
		DevReset: DevReset,
		Migrate:  Migrate,
		AddSeed:  AddSeed,
	})
}

func Close() error {
	if inst == nil {
		return nil
	}
	sqlDB, err := inst.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
