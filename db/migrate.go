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
