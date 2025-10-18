package db

import (
	"uszatakuchnia/db/entities"

	"gorm.io/gorm"
)

func AutoMigrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		&entities.TasteType{},
		&entities.IngredientType{},
		&entities.AromaType{},
		&entities.Ingredient{},
		&entities.Taste{},
		&entities.IngredientTaste{},
		&entities.Aroma{},
	)
}

func runMigrations(gdb *gorm.DB) error {
	_ = gdb.Exec(`SET LOCAL statement_timeout = '30s'`).Error

	const lockKey int64 = 2137
	if err := gdb.Exec(`SELECT pg_advisory_lock(?)`, lockKey).Error; err != nil {
		return err
	}

	defer gdb.Exec(`SELECT pg_advisory_unlock(?)`, lockKey)

	if err := AutoMigrate(gdb); err != nil {
		return err
	}

	return nil
}

func devReset(gdb *gorm.DB) error {
	return gdb.Exec(`DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;`).Error
}
