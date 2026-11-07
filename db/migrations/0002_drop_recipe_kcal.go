package migrations

import "gorm.io/gorm"

// 0002: ręczne kcal_per_serving zastąpione wyliczanym recipes.nutrition (makro ze składników).
var dropRecipeKcal = Migration{
	ID: "0002_drop_recipe_kcal",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasColumn("recipes", "kcal_per_serving") {
			return nil // świeża baza
		}
		return tx.Migrator().DropColumn("recipes", "kcal_per_serving")
	},
}
