package db

import (
	"sync"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/migrations"

	"gorm.io/gorm"
)

func AutoMigrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		&entities.IngredientType{},
		&entities.Ingredient{},
		&entities.IngredientAllergen{},
		&entities.IngredientUnitConversion{},

		&entities.TaxonomyAxis{},
		&entities.TaxonomyValue{},
		&entities.Unit{},

		&entities.Recipe{},
		&entities.RecipeStep{},
		&entities.RecipeIngredient{},
		&entities.RecipePhoto{},
		&entities.RecipeTaxonomy{},

		&entities.Image{},
		&entities.Photo{},
	)
}

// runMigrations: najpierw AutoMigrate (dodaje kolumny/tabele), potem wersjonowane migracje danych.
func runMigrations(gdb *gorm.DB) error {
	if err := AutoMigrate(gdb); err != nil {
		return err
	}
	return migrations.Run(gdb)
}

var migrateLock sync.Mutex

// devResetTables — wszystkie tabele aplikacji, łącznie z wycofanymi (żeby reset nie zostawiał śmieci).
var devResetTables = []string{
	"recipe_taxonomy",
	"recipe_tags", // wycofana: zastąpiona przez recipe_taxonomy
	"recipe_steps",
	"recipe_ingredients",
	"recipe_photos",
	"images",
	"photos",
	"recipes",
	"taxonomy_values",
	"taxonomy_axes",
	"units",
	"ingredient_allergens",
	"ingredients",
	"ingredient_types",
	"schema_migrations",
}

func devReset(db *gorm.DB) error {
	migrateLock.Lock()
	defer migrateLock.Unlock()

	return db.Transaction(func(tx *gorm.DB) error {
		for _, t := range devResetTables {
			// postgres migrator: DROP TABLE IF EXISTS ... CASCADE
			if err := tx.Migrator().DropTable(t); err != nil {
				return err
			}
		}
		return nil
	})
}
