package entities

// RecipeIngredientSubstitute — zamiennik składnika w wierszu przepisu (recipe_ingredients: recipe_id + sort_order).
// Mix = „pół na pół”: składnik główny + zamiennik, po połowie ilości. Zamiennik liczy się jako użycie składnika
// (ON DELETE RESTRICT). Usunięcie wiersza przepisu usuwa jego zamienniki (FK złożony, migracja 0008).
type RecipeIngredientSubstitute struct {
	RecipeID  uint `gorm:"not null;primaryKey;uniqueIndex:ux_recsub_option,priority:1"`
	SortOrder int  `gorm:"not null;primaryKey;uniqueIndex:ux_recsub_option,priority:2"`
	Position  int  `gorm:"not null;primaryKey"`

	IngredientID uint       `gorm:"not null;index:ix_recsub_ingredient;uniqueIndex:ux_recsub_option,priority:3"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Mix bool `gorm:"not null;default:false;uniqueIndex:ux_recsub_option,priority:4"`
}

func (RecipeIngredientSubstitute) TableName() string { return "recipe_ingredient_substitutes" }
