package db

import (
	"sync"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/migrations"

	"gorm.io/gorm"
)

// models lists every entity with a table (AutoMigrate and the schema fingerprint in schema_state.go).
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

	&entities.ChallengeFormat{},
	&entities.Challenge{},
	&entities.ChallengeIngredient{},

	&entities.GlossaryTerm{},

	&entities.Image{},
	&entities.Photo{},
}

func AutoMigrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(models...)
}

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

// devResetTables lists every app table, retired ones included, so a reset leaves no leftovers.
var devResetTables = []string{
	"glossary_terms",
	"challenge_ingredients",
	"challenges",
	"challenge_formats",
	"recipe_equipment",
	"equipment",
	"recipe_taxonomy",
	"recipe_tags", // retired: replaced by recipe_taxonomy
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
