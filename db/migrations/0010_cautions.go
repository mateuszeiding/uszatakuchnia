package migrations

import "gorm.io/gorm"

// 0010: recipes.cautions i challenges.cautions — uwagi „czego nie zawalić” (kolumny zakłada AutoMigrate,
// jsonb NOT NULL DEFAULT '[]'). Zabezpieczenie: wiersze bez wartości dostają pustą listę.
var cautions = Migration{
	ID: "0010_cautions",
	Up: func(tx *gorm.DB) error {
		for _, t := range []string{"recipes", "challenges"} {
			if !tx.Migrator().HasColumn(t, "cautions") {
				continue
			}
			if err := tx.Exec(`UPDATE ` + t + ` SET cautions = '[]'::jsonb WHERE cautions IS NULL`).Error; err != nil {
				return err
			}
		}
		return nil
	},
}
