package entities

type RecipeIngredient struct {
	RecipeID  uint `gorm:"not null;primaryKey"`
	SortOrder int  `gorm:"not null;primaryKey"`

	Section *string `gorm:"type:varchar(64)"`

	IngredientID uint       `gorm:"not null;index:ix_recing_ingredient"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Amount     *float64 `gorm:"type:decimal(10,3)"`
	Unit       *string  `gorm:"type:varchar(32)"`
	AmountText *string  `gorm:"type:varchar(255)"`
	Note       *string  `gorm:"type:varchar(255)"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeIngredient) TableName() string { return "recipe_ingredients" }
