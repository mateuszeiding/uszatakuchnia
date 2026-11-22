package entities

// Portion modes: what one "serving" of the recipe means to the reader.
const (
	PortionServings = "servings" // people / pieces, the classic counter
	PortionSized    = "sized"    // pieces in author-defined sizes (pizza): servings = pieces in the recipe
	PortionForm     = "form"     // one baking form or tray, scaled by the area of its base
)

const (
	FormRound = "round"
	FormRect  = "rect"
)

// Form and piece sizes are whole centimetres in this range (the reader's own form too).
const (
	PortionMinCm = 10
	PortionMaxCm = 60
)

// RecipePieceSize is one size the reader can pick in "sized" mode. Exactly one per published recipe is the
// base (the size the ingredient quantities are written for); ux_recipe_piece_sizes_base (migration 0014)
// keeps it at most one.
type RecipePieceSize struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	RecipeID   uint `gorm:"not null;uniqueIndex:ux_recipe_piece_sizes_diameter,priority:1"`
	DiameterCm int  `gorm:"not null;uniqueIndex:ux_recipe_piece_sizes_diameter,priority:2"`
	DoughG     int  `gorm:"not null"`
	IsBase     bool `gorm:"not null;default:false"`
	SortOrder  int  `gorm:"not null;default:0"`
}

func (RecipePieceSize) TableName() string { return "recipe_piece_sizes" }

// NutritionDivisor is what the whole-recipe macros are divided by for recipes.nutrition.perPortion:
// one serving, one piece of the base size, or one piece of the form (the whole form without pieces).
func (r Recipe) NutritionDivisor() int {
	if r.PortionMode == PortionForm {
		if r.FormPieces != nil && *r.FormPieces > 0 {
			return *r.FormPieces
		}
		return 1
	}
	return r.Servings
}
