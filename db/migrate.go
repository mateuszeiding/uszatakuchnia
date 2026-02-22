package db

import (
	"sync"
	"uszatakuchnia/db/entities"

	"gorm.io/gorm"
)

func AutoMigrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		&entities.IngredientType{},
		&entities.Ingredient{},

		&entities.Recipe{},
		&entities.RecipeStep{},
		&entities.RecipeIngredient{},
		&entities.RecipePhoto{},

		&entities.Image{},
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
		// Dropuj w odwrotnej kolejności zależności (children -> parents)
		if err := tx.Migrator().DropTable(
			&entities.RecipeStep{},       // zależy od recipes
			&entities.RecipeIngredient{}, // zależy od recipes i ingredients
			&entities.Image{},            // polimorficzne - zależy logicznie od encji, ale bez FK
			&entities.RecipePhoto{},

			&entities.Recipe{},     // parent dla steps/ingredients
			&entities.Ingredient{}, // parent dla recipe_ingredients

			&entities.IngredientType{}, // parent dla ingredients
		); err != nil {
			return err
		}
		return nil
	})
}
