package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// 0012: the glossary (glossary_terms). AutoMigrate creates the table; here the name becomes unique case-insensitively
// (like equipment in 0007). Moves no data: the table is new and empty.
var glossaryTerms = Migration{
	ID: "0012_glossary_terms",
	Up: func(tx *gorm.DB) error {
		if !tx.Migrator().HasTable("glossary_terms") {
			return fmt.Errorf("table glossary_terms missing after AutoMigrate")
		}
		return tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS ux_glossary_terms_name_lower ON glossary_terms (lower(name))`).Error
	},
}
