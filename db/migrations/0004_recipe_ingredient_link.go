package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// 0004: recipe_ingredients.linked_recipe_id (link ze składnika do przepisu). Kolumnę i FK zakłada AutoMigrate;
// tu pilnujemy, żeby usunięcie przepisu zerowało link (ON DELETE SET NULL), a nie blokowało usunięcia.
var recipeIngredientLink = Migration{
	ID: "0004_recipe_ingredient_link",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasColumn("recipe_ingredients", "linked_recipe_id") {
			return fmt.Errorf("brak kolumny recipe_ingredients.linked_recipe_id po AutoMigrate")
		}
		var fks []struct {
			Name string
			Rule string
		}
		if err := tx.Raw(`
			SELECT rc.constraint_name AS name, rc.delete_rule AS rule
			FROM information_schema.referential_constraints rc
			JOIN information_schema.key_column_usage k
			  ON k.constraint_name = rc.constraint_name AND k.constraint_schema = rc.constraint_schema
			WHERE k.table_name = 'recipe_ingredients' AND k.column_name = 'linked_recipe_id'`).
			Scan(&fks).Error; err != nil {
			return err
		}
		for _, fk := range fks {
			if fk.Rule == "SET NULL" {
				return nil
			}
		}
		for _, fk := range fks {
			if err := tx.Exec(fmt.Sprintf(`ALTER TABLE recipe_ingredients DROP CONSTRAINT %q`, fk.Name)).Error; err != nil {
				return err
			}
		}
		return tx.Exec(`ALTER TABLE recipe_ingredients
			ADD CONSTRAINT fk_recipe_ingredients_linked_recipe FOREIGN KEY (linked_recipe_id)
			REFERENCES recipes(id) ON UPDATE CASCADE ON DELETE SET NULL`).Error
	},
}
