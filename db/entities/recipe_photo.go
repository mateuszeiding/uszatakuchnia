package entities

import dtos "uszatakuchnia/dtos"

type RecipePhoto struct {
	RecipeID uint   `gorm:"primaryKey;not null"`
	ImageURL string `gorm:"not null"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipePhoto) TableName() string { return "recipe_photos" }

func (p *RecipePhoto) ToDto() *dtos.RecipePhotoDto {
	if p == nil {
		return nil
	}

	return &dtos.RecipePhotoDto{
		URL: p.ImageURL,
	}
}
