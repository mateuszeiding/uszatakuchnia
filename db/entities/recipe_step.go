package entities

type RecipeStep struct {
	RecipeID uint    `gorm:"not null;primaryKey"`
	StepNo   int     `gorm:"not null;primaryKey"`
	Section  *string `gorm:"type:varchar(64)"`
	Text     string  `gorm:"not null"`

	// Step to do ahead (marinade, soaking). A recipe "needs prep" = it has such a step that isn't optional.
	Prep     bool    `gorm:"not null;default:false"`
	PrepLead *string `gorm:"type:varchar(32)"`

	// A step the cook may skip: progress and "needs prep" count only the other steps.
	Optional bool `gorm:"not null;default:false"`

	Hints StringList `gorm:"type:jsonb;not null;default:'[]'"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeStep) TableName() string { return "recipe_steps" }
