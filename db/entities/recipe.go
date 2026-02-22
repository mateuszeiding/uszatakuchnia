package entities

import dtos "uszatakuchnia/dtos"

type Recipe struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);not null;index:ix_recipes_name"`
	Servings    int    `gorm:"not null"` // CHECK > 0 ogarnij migracją
	Description *string

	Photo       *RecipePhoto       `gorm:"foreignKey:RecipeID;references:ID"`
	Steps       []RecipeStep       `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Ingredients []RecipeIngredient `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Recipe) TableName() string { return "recipes" }

func (r Recipe) ToDto() dtos.RecipeDto {
	return dtos.RecipeDto{
		ID:          r.ID,
		Name:        r.Name,
		Servings:    r.Servings,
		Description: r.Description,
		Photo:       r.Photo.ToDto(),
	}
}
