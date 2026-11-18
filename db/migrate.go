package db

import (
	"sync"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/migrations"

	"gorm.io/gorm"
)

// models — wszystkie encje z tabelami (AutoMigrate i odcisk schematu w schema_state.go).
var models = []any{
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
	&entities.RecipeIngredientSubstitute{},
	&entities.RecipePhoto{},
	&entities.RecipeTaxonomy{},

	&entities.Equipment{},
	&entities.RecipeEquipment{},

	&entities.Image{},
	&entities.Photo{},
}

func AutoMigrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(models...)
}

// runMigrations: najpierw AutoMigrate (dodaje kolumny/tabele), potem wersjonowane migracje danych.
func runMigrations(gdb *gorm.DB) error {
	if err := AutoMigrate(gdb); err != nil {
		return err
	}
	if err := migrations.Run(gdb); err != nil {
		return err
	}
	return saveFingerprint(gdb)
}

var migrateLock sync.Mutex

// devResetTables — wszystkie tabele aplikacji, łącznie z wycofanymi (żeby reset nie zostawiał śmieci).
var devResetTables = []string{
	"recipe_equipment",
	"equipment",
	"recipe_taxonomy",
	"recipe_tags", // wycofana: zastąpiona przez recipe_taxonomy
	"recipe_steps",
	"recipe_ingredient_substitutes",
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
	"schema_state",
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
