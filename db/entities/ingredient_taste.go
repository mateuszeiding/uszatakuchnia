package entities

type IngredientTaste struct {
	IngredientID int `gorm:"primaryKey;not null"`
	TasteID      int `gorm:"primaryKey;not null;index:ix_ingtaste_taste"`
	Intensity    int `gorm:"not null;default:50"`

	Ingredient Ingredient `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Taste      Taste      `gorm:"foreignKey:TasteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (IngredientTaste) TableName() string { return "ingredient_taste" }
