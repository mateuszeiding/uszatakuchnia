package entities

import "time"

// GlossaryTerm is a cooking term explained in steps and hints ({{term:<id>|<word>}} in their text).
// Name is unique case-insensitively (index ux_glossary_terms_name_lower, migration 0012). Which recipes use
// a term is read from the step text, so there is no join table that could drift from it.
type GlossaryTerm struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(60);not null"`
	Explanation string    `gorm:"type:varchar(280);not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (GlossaryTerm) TableName() string { return "glossary_terms" }
