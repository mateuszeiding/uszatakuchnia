package entities

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
)

type Recipe struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);not null;index:ix_recipes_name"`
	Servings    int    `gorm:"not null;default:0"` // 0 = not set (draft)
	Description *string
	Tagline     *string `gorm:"type:varchar(500)"`
	Tip         *string
	// "Don't mess this up" notes (shopping, safety, timing) in the author's order; nil is stored as [].
	Cautions StringList `gorm:"type:jsonb;not null;default:'[]'"`

	// Time in minutes, from–to. TimeMin 0 = not given (draft); TimeMax nil = no range.
	TimeMin    int `gorm:"not null;default:0"`
	TimeMax    *int
	Difficulty *int
	// Macro cache (services.RecomputeRecipes); nil = no data.
	Nutrition *RecipeNutrition `gorm:"type:jsonb"`

	Status string `gorm:"type:varchar(16);not null;default:'draft'"`

	Photo       *RecipePhoto       `gorm:"foreignKey:RecipeID;references:ID"`
	Steps       []RecipeStep       `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Ingredients []RecipeIngredient `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Taxonomy    []RecipeTaxonomy   `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Equipment   []RecipeEquipment  `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// RecipeBase is a light recipe view for the list. NeedsPrep is computed in the query (EXISTS on prep steps).
type RecipeBase struct {
	ID         uint    `gorm:"primaryKey;autoIncrement"`
	Name       string  `gorm:"type:varchar(255);not null;index:ix_recipes_name"`
	Tagline    *string `gorm:"type:varchar(500)"`
	TimeMin    int
	TimeMax    *int
	Difficulty *int
	Status     string           `gorm:"type:varchar(16);not null"`
	NeedsPrep  bool             `gorm:"->;-:migration"`
	Nutrition  *RecipeNutrition `gorm:"->;-:migration"`
	Photo      *RecipePhoto     `gorm:"foreignKey:RecipeID;references:ID"`
	Taxonomy   []RecipeTaxonomy `gorm:"foreignKey:RecipeID;references:ID"`
}

func (Recipe) TableName() string     { return "recipes" }
func (RecipeBase) TableName() string { return "recipes" }
