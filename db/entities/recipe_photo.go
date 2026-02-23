package entities

type RecipePhoto struct {
	RecipeID uint   `gorm:"primaryKey;not null"`
	ImageURL string `gorm:"not null"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipePhoto) TableName() string { return "recipe_photos" }
