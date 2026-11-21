package migrations

import "gorm.io/gorm"

// 0011: recipe_ingredients.key — a stable row key that step text refers to ({{ing:<key>}}). AutoMigrate adds the
// column with an empty key for existing rows; here they get "i<sort_order>" (unique within a recipe), and only
// then the unique index can be built. The index skips empty keys, so a row written without one (restore-recipe
// from an older backup, before it fills them in) can't break the insert.
var recipeIngredientKey = Migration{
	ID: "0011_recipe_ingredient_key",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasColumn("recipe_ingredients", "key") {
			return nil
		}
		if err := tx.Exec(`UPDATE recipe_ingredients SET "key" = 'i' || sort_order WHERE "key" = ''`).Error; err != nil {
			return err
		}
		return tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS ux_recing_key ON recipe_ingredients (recipe_id, "key") WHERE "key" <> ''`).Error
	},
}
