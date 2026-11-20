package entities

// Mix = "half and half": main ingredient + substitute, half the quantity each. A substitute counts as a use of the
// ingredient (ON DELETE RESTRICT). Deleting the recipe row deletes its substitutes (composite FK, migration 0008).
type RecipeIngredientSubstitute struct {
	RecipeID  uint `gorm:"not null;primaryKey;uniqueIndex:ux_recsub_option,priority:1"`
	SortOrder int  `gorm:"not null;primaryKey;uniqueIndex:ux_recsub_option,priority:2"`
	Position  int  `gorm:"not null;primaryKey"`

	IngredientID uint       `gorm:"not null;index:ix_recsub_ingredient;uniqueIndex:ux_recsub_option,priority:3"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Mix bool `gorm:"not null;default:false;uniqueIndex:ux_recsub_option,priority:4"`
}

func (RecipeIngredientSubstitute) TableName() string { return "recipe_ingredient_substitutes" }
