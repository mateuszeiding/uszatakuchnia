package entities

import "time"

// ChallengeFormat — słownik formatów wyzwań („z podanych składników”, „na temat”, „odtwórz danie”…).
// Nazwa unikalna bez wielkości liter (indeks ux_challenge_formats_name_lower, migracja 0009).
type ChallengeFormat struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(64);not null"`
	SortOrder int    `gorm:"not null;default:0"`
}

func (ChallengeFormat) TableName() string { return "challenge_formats" }

// Challenge — wyzwanie kulinarne. Bez statusu (od razu publiczne) i bez powiązań z przepisami.
// CHECK: difficulty 1–5, time_limit_min > 0 (migracja 0009).
type Challenge struct {
	ID          uint            `gorm:"primaryKey;autoIncrement"`
	Title       string          `gorm:"type:varchar(255);not null"`
	FormatID    uint            `gorm:"not null;index:ix_challenges_format"`
	Format      ChallengeFormat `gorm:"foreignKey:FormatID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Description string          `gorm:"type:text;not null"` // akapity rozdzielone pustą linią
	// nil = bez limitu. Tylko informacja, bez zegara.
	TimeLimitMin *int
	Difficulty   int     `gorm:"not null"` // 1 rozgrzewka … 5 mistrzowskie
	Source       *string `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// Uwagi „czego nie zawalić” w kolejności autora; nil zapisuje się jako [].
	Cautions StringList `gorm:"type:jsonb;not null;default:'[]'"`

	Ingredients []ChallengeIngredient `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Challenge) TableName() string { return "challenges" }

// ChallengeIngredient — pozycja tekstowa („ryba białego mięsa”), bez powiązania z bazą składników.
type ChallengeIngredient struct {
	ChallengeID uint   `gorm:"not null;primaryKey"`
	SortOrder   int    `gorm:"not null;primaryKey"`
	Name        string `gorm:"type:varchar(120);not null"`
	Required    bool   `gorm:"not null;default:false"` // „Musisz użyć”
}

func (ChallengeIngredient) TableName() string { return "challenge_ingredients" }
