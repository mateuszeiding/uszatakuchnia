package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// 0009: wyzwania. Tabele i FK zakłada AutoMigrate (challenge_ingredients ON DELETE CASCADE,
// format_id ON DELETE RESTRICT); tu ograniczenia CHECK, unikalna nazwa formatu bez wielkości liter
// i formaty startowe (na świeżej bazie też).
var challenges = Migration{
	ID: "0009_challenges",
	Up: func(tx *gorm.DB) error {
		for _, t := range []string{"challenge_formats", "challenges", "challenge_ingredients"} {
			if !tx.Migrator().HasTable(t) {
				return fmt.Errorf("brak tabeli %s po AutoMigrate", t)
			}
		}
		stmts := []string{
			`CREATE UNIQUE INDEX IF NOT EXISTS ux_challenge_formats_name_lower ON challenge_formats (lower(name))`,
			`ALTER TABLE challenges DROP CONSTRAINT IF EXISTS ck_challenges_difficulty`,
			`ALTER TABLE challenges ADD CONSTRAINT ck_challenges_difficulty CHECK (difficulty BETWEEN 1 AND 5)`,
			`ALTER TABLE challenges DROP CONSTRAINT IF EXISTS ck_challenges_time_limit`,
			`ALTER TABLE challenges ADD CONSTRAINT ck_challenges_time_limit CHECK (time_limit_min IS NULL OR time_limit_min > 0)`,
		}
		for _, s := range stmts {
			if err := tx.Exec(s).Error; err != nil {
				return err
			}
		}
		for i, name := range []string{"z podanych składników", "na temat", "odtwórz danie"} {
			if err := tx.Exec(
				`INSERT INTO challenge_formats (name, sort_order)
				 SELECT ?, ? WHERE NOT EXISTS (SELECT 1 FROM challenge_formats WHERE lower(name) = lower(?))`,
				name, i+1, name,
			).Error; err != nil {
				return err
			}
		}
		return nil
	},
}
