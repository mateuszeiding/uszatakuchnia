package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// 0007: słownik sprzętu (equipment) i sprzęt przepisu (recipe_equipment). Tabele i FK zakłada AutoMigrate
// (recipe_id ON DELETE CASCADE, equipment_id ON DELETE RESTRICT); tu unikalna nazwa bez wielkości liter.
var equipment = Migration{
	ID: "0007_equipment",
	Up: func(tx *gorm.DB) error {
		for _, t := range []string{"equipment", "recipe_equipment"} {
			if !tx.Migrator().HasTable(t) {
				return fmt.Errorf("brak tabeli %s po AutoMigrate", t)
			}
		}
		return tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS ux_equipment_name_lower ON equipment (lower(name))`).Error
	},
}
