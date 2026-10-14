package migrations

import "gorm.io/gorm"

// 0001: czas przepisu jako zakres. recipes.time_minutes (nullable) → time_min (NOT NULL, 0 = nie podano).
// Kolumny time_min / time_max dodaje AutoMigrate; tu tylko przenosimy dane i usuwamy starą kolumnę.
var recipeTimeRange = Migration{
	ID: "0001_recipe_time_range",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasColumn("recipes", "time_minutes") {
			return nil // świeża baza
		}
		if err := tx.Exec("UPDATE recipes SET time_min = time_minutes WHERE time_minutes IS NOT NULL").Error; err != nil {
			return err
		}
		return tx.Migrator().DropColumn("recipes", "time_minutes")
	},
}
