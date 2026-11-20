package entities

import "time"

// ChallengeFormat is the dictionary of challenge formats ("z podanych składników", "na temat", "odtwórz danie"…).
// Name is unique case-insensitively (index ux_challenge_formats_name_lower, migration 0009).
type ChallengeFormat struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(64);not null"`
	SortOrder int    `gorm:"not null;default:0"`
}

func (ChallengeFormat) TableName() string { return "challenge_formats" }

// Challenge is a cooking challenge. No status (public right away) and no links to recipes.
// CHECK: difficulty 1–5, time_limit_min > 0 (migration 0009).
type Challenge struct {
	ID          uint            `gorm:"primaryKey;autoIncrement"`
	Title       string          `gorm:"type:varchar(255);not null"`
	FormatID    uint            `gorm:"not null;index:ix_challenges_format"`
	Format      ChallengeFormat `gorm:"foreignKey:FormatID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Description string          `gorm:"type:text;not null"` // paragraphs separated by a blank line
	// nil = no limit. Informational only, no timer.
	TimeLimitMin *int
	Difficulty   int     `gorm:"not null"` // 1 warm-up … 5 masterclass
	Source       *string `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// "Don't mess this up" notes in the author's order; nil is stored as [].
	Cautions StringList `gorm:"type:jsonb;not null;default:'[]'"`

	Ingredients []ChallengeIngredient `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Challenge) TableName() string { return "challenges" }

// ChallengeIngredient is a free-text item ("ryba białego mięsa"), not linked to the ingredients table.
type ChallengeIngredient struct {
	ChallengeID uint   `gorm:"not null;primaryKey"`
	SortOrder   int    `gorm:"not null;primaryKey"`
	Name        string `gorm:"type:varchar(120);not null"`
	Required    bool   `gorm:"not null;default:false"` // "Musisz użyć" (must use)
}

func (ChallengeIngredient) TableName() string { return "challenge_ingredients" }
