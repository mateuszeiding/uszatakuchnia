package migrations

import "gorm.io/gorm"

// 0006: recipe_steps.hints (kolumnę zakłada AutoMigrate, jsonb NOT NULL DEFAULT '[]') — wskazówki przy krokach.
// Zabezpieczenie: wiersze bez wartości dostają pustą listę. recipes.tip zostaje do czasu przeniesienia w formularzu.
var recipeStepHints = Migration{
	ID: "0006_recipe_step_hints",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasColumn("recipe_steps", "hints") {
			return nil
		}
		return tx.Exec(`UPDATE recipe_steps SET hints = '[]'::jsonb WHERE hints IS NULL`).Error
	},
}
