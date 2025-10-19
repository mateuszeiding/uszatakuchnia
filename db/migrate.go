package db

import (
	"sync"
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

var migrateLock sync.Mutex

func devReset(db *gorm.DB) error {
	migrateLock.Lock()
	defer migrateLock.Unlock()

	return db.Transaction(func(tx *gorm.DB) error {
		// Dropuj w odwrotnej kolejności zależności
		if err := tx.Migrator().DropTable(
			&entities.IngredientTaste{}, // PK (ingredient_id, taste_id)
			&entities.Aroma{},
			&entities.Taste{},
			&entities.Ingredient{},
			&entities.TasteType{},
			&entities.IngredientType{},
			&entities.AromaType{},
		); err != nil {
			return err
		}
		return nil
	})
}
