package entities

type RecipeStep struct {
	RecipeID uint    `gorm:"not null;primaryKey"`
	StepNo   int     `gorm:"not null;primaryKey"`
	Section  *string `gorm:"type:varchar(64)"`
	Text     string  `gorm:"not null"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeStep) TableName() string { return "recipe_steps" }
