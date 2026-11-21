package migrations

import "gorm.io/gorm"

// 0013: taxonomy_values.icon — the dish-type icon on the recipe card placeholder frame (the column comes
// from AutoMigrate). Sets the icon of the nine dish types that exist today, matched by label on the "dish"
// axis, only where it is still empty. Renamed or new types are left alone (no icon = color only).
// On a fresh database the table is empty here; the seed sets the icons itself.
var taxonomyValueIcons = Migration{
	ID: "0013_taxonomy_value_icons",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasColumn("taxonomy_values", "icon") {
			return nil
		}
		byLabel := []struct{ label, icon string }{
			{"śniadanie", "sniadanie"},
			{"zupa", "zupa"},
			{"danie główne", "glowne"},
			{"sałatka", "salatka"},
			{"przekąska", "przekaska"},
			{"deser", "deser"},
			{"pieczywo i wypieki", "wypiek"},
			{"sos / dodatek", "sos"},
			{"napój", "napoj"},
		}
		for _, v := range byLabel {
			err := tx.Exec(`UPDATE taxonomy_values SET icon = ?
				WHERE icon IS NULL AND lower(label) = lower(?)
				AND axis_id IN (SELECT id FROM taxonomy_axes WHERE code = 'dish')`, v.icon, v.label).Error
			if err != nil {
				return err
			}
		}
		return nil
	},
}
