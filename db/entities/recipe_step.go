package entities

import dtos "uszatakuchnia/dtos"

type RecipeStep struct {
	RecipeID uint    `gorm:"not null;primaryKey"`
	StepNo   int     `gorm:"not null;primaryKey"`
	Section  *string `gorm:"type:varchar(64)"`
	Text     string  `gorm:"not null"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeStep) TableName() string { return "recipe_steps" }

func (s RecipeStep) ToDto() dtos.RecipeStepDto {
	return dtos.RecipeStepDto{
		StepNo:  s.StepNo,
		Section: s.Section,
		Text:    s.Text,
	}
}
