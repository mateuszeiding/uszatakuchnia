package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// 0008: zamienniki składników (recipe_ingredient_substitutes). Tabelę, FK do ingredients (RESTRICT), FK złożony
// do wiersza przepisu i unikalność (recipe_id, sort_order, ingredient_id, mix) zakłada AutoMigrate.
// Tu pilnujemy, żeby FK do recipe_ingredients kasował zamienniki razem z wierszem (zapis formularza zastępuje
// wiersze) — AutoMigrate nie poprawia reguły istniejącego FK.
var recipeIngredientSubstitutes = Migration{
	ID: "0008_recipe_ingredient_substitutes",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasTable("recipe_ingredient_substitutes") {
			return fmt.Errorf("brak tabeli recipe_ingredient_substitutes po AutoMigrate")
		}
		var fks []struct {
			Name string
			Rule string
		}
		if err := tx.Raw(`
			SELECT DISTINCT rc.constraint_name AS name, rc.delete_rule AS rule
			FROM information_schema.referential_constraints rc
			JOIN information_schema.table_constraints t
			  ON t.constraint_name = rc.constraint_name AND t.constraint_schema = rc.constraint_schema
			JOIN information_schema.constraint_column_usage u
			  ON u.constraint_name = rc.constraint_name AND u.constraint_schema = rc.constraint_schema
			WHERE t.table_name = 'recipe_ingredient_substitutes' AND u.table_name = 'recipe_ingredients'`).
			Scan(&fks).Error; err != nil {
			return err
		}
		cascade := false
		for _, fk := range fks {
			if fk.Rule == "CASCADE" {
				cascade = true
				continue
			}
			if err := tx.Exec(fmt.Sprintf(`ALTER TABLE recipe_ingredient_substitutes DROP CONSTRAINT %q`, fk.Name)).Error; err != nil {
				return err
			}
		}
		if cascade {
			return nil
		}
		return tx.Exec(`ALTER TABLE recipe_ingredient_substitutes
			ADD CONSTRAINT fk_recipe_ingredients_substitutes FOREIGN KEY (recipe_id, sort_order)
			REFERENCES recipe_ingredients(recipe_id, sort_order) ON UPDATE CASCADE ON DELETE CASCADE`).Error
	},
}
