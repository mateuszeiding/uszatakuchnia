package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// 0014: portion size (recipes.portion_mode, piece unit, form fields, recipe_piece_sizes — all from AutoMigrate;
// portion_mode defaults to 'servings', so existing recipes don't change).
// Here: at most one base size per recipe, and the "placek" unit for pizza-like recipes. The unit is added only
// to a database that already has units: on a fresh one the table is still empty here and the seed adds it in
// its own order.
var portionSize = Migration{
	ID: "0014_portion_size",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasTable("recipe_piece_sizes") {
			return fmt.Errorf("table recipe_piece_sizes missing after AutoMigrate")
		}
		if err := tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS ux_recipe_piece_sizes_base
			ON recipe_piece_sizes (recipe_id) WHERE is_base`).Error; err != nil {
			return err
		}
		var units int64
		if err := tx.Table("units").Count(&units).Error; err != nil {
			return err
		}
		if units == 0 {
			return nil
		}
		return tx.Exec(`INSERT INTO units (name, form_one, form_few, form_many, no_qty, sort_order)
			SELECT 'placek', 'placek', 'placki', 'placków', false, COALESCE(MAX(sort_order), 0) + 1 FROM units
			ON CONFLICT (name) DO NOTHING`).Error
	},
}
