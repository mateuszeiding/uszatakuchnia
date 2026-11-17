package migrations

import "gorm.io/gorm"

// 0005: ingredients.skip_shopping (kolumnę zakłada AutoMigrate, domyślnie false) — woda nie trafia na listę zakupów.
var ingredientSkipShopping = Migration{
	ID: "0005_ingredient_skip_shopping",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasColumn("ingredients", "skip_shopping") {
			return nil
		}
		return tx.Exec(`UPDATE ingredients SET skip_shopping = true WHERE lower(trim(name)) = 'woda'`).Error
	},
}
