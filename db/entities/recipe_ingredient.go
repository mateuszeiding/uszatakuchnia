package entities

type RecipeIngredient struct {
	RecipeID  uint `gorm:"not null;primaryKey"`
	SortOrder int  `gorm:"not null;primaryKey"`

	Section *string `gorm:"type:varchar(64)"`

	IngredientID uint       `gorm:"not null;index:ix_recing_ingredient"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Ilość jako tekst: liczba, ułamek ("½", "1/4", "1,5") lub zakres ("2–3"). Parsowana przez services.ParseQuantity.
	Quantity *string `gorm:"type:varchar(32)"`
	UnitID   *uint   `gorm:"index:ix_recing_unit"`
	Unit     *Unit   `gorm:"foreignKey:UnitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Note     *string `gorm:"type:varchar(255)"`
	Optional bool    `gorm:"not null;default:false"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeIngredient) TableName() string { return "recipe_ingredients" }
