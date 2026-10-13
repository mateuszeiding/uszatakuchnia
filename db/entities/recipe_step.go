package entities

type RecipeStep struct {
	RecipeID uint    `gorm:"not null;primaryKey"`
	StepNo   int     `gorm:"not null;primaryKey"`
	Section  *string `gorm:"type:varchar(64)"`
	Text     string  `gorm:"not null"`

	// Krok do zrobienia wcześniej (marynata, namaczanie). Przepis "wymaga przygotowania" = ma taki krok.
	Prep     bool    `gorm:"not null;default:false"`
	PrepLead *string `gorm:"type:varchar(32)"` // np. "12 h", "dzień wcześniej"

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeStep) TableName() string { return "recipe_steps" }
