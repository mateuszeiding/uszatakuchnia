package entities

type RecipeStep struct {
	RecipeID uint    `gorm:"not null;primaryKey"`
	StepNo   int     `gorm:"not null;primaryKey"`
	Section  *string `gorm:"type:varchar(64)"`
	Text     string  `gorm:"not null"`

	// Step to do ahead (marinade, soaking). A recipe "needs prep" = it has such a step.
	Prep     bool    `gorm:"not null;default:false"`
	PrepLead *string `gorm:"type:varchar(32)"` // e.g. "12 h", "dzień wcześniej"

	// Hints below the step text ("jeśli…", "możesz też…"), 0–n.
	Hints StringList `gorm:"type:jsonb;not null;default:'[]'"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeStep) TableName() string { return "recipe_steps" }
