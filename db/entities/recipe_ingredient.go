package entities

type RecipeIngredient struct {
	RecipeID  uint `gorm:"not null;primaryKey"`
	SortOrder int  `gorm:"not null;primaryKey"`

	Section *string `gorm:"type:varchar(64)"`

	// Stable row key within the recipe: step text refers to the row by it ({{ing:<key>}}), and sort_order
	// changes whenever rows are reordered. The unique index (recipe_id, key) lives in migration 0011, because
	// AutoMigrate would try to build it over the empty defaults of existing rows.
	Key string `gorm:"type:varchar(32);not null;default:''"`

	IngredientID uint       `gorm:"not null;index:ix_recing_ingredient"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Quantity as text: a number, a fraction ("½", "1/4", "1,5") or a range ("2–3"). Parsed by services.ParseQuantity.
	Quantity *string `gorm:"type:varchar(32)"`
	UnitID   *uint   `gorm:"index:ix_recing_unit"`
	Unit     *Unit   `gorm:"foreignKey:UnitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Note     *string `gorm:"type:varchar(255)"`
	Optional bool    `gorm:"not null;default:false"`

	// Deleting the linked recipe clears the link.
	LinkedRecipeID *uint   `gorm:"index:ix_recing_linked_recipe"`
	LinkedRecipe   *Recipe `gorm:"foreignKey:LinkedRecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// Substitutes ordered by position. Deleting the row (saving the form replaces rows) deletes its substitutes.
	Substitutes []RecipeIngredientSubstitute `gorm:"foreignKey:RecipeID,SortOrder;references:RecipeID,SortOrder;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeIngredient) TableName() string { return "recipe_ingredients" }
