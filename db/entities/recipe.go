package entities

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
)

type Recipe struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);not null;index:ix_recipes_name"`
	Servings    int    `gorm:"not null;default:0"` // 0 = nieustawione (szkic)
	Description *string
	Tagline     *string `gorm:"type:varchar(500)"`
	Tip         *string // wskazówka pod krokami

	// Czas w minutach: od–do. TimeMin 0 = nie podano (szkic); TimeMax nil = bez zakresu.
	TimeMin        int `gorm:"not null;default:0"`
	TimeMax        *int
	Difficulty     *int
	KcalPerServing *int

	Status string `gorm:"type:varchar(16);not null;default:'draft'"`

	Photo       *RecipePhoto       `gorm:"foreignKey:RecipeID;references:ID"`
	Steps       []RecipeStep       `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Ingredients []RecipeIngredient `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Taxonomy    []RecipeTaxonomy   `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// RecipeBase — lekki widok przepisu dla listy. NeedsPrep wyliczane w zapytaniu (EXISTS na krokach prep).
type RecipeBase struct {
	ID         uint    `gorm:"primaryKey;autoIncrement"`
	Name       string  `gorm:"type:varchar(255);not null;index:ix_recipes_name"`
	Tagline    *string `gorm:"type:varchar(500)"`
	TimeMin    int
	TimeMax    *int
	Difficulty *int
	Status     string           `gorm:"type:varchar(16);not null"`
	NeedsPrep  bool             `gorm:"->;-:migration"`
	Photo      *RecipePhoto     `gorm:"foreignKey:RecipeID;references:ID"`
	Taxonomy   []RecipeTaxonomy `gorm:"foreignKey:RecipeID;references:ID"`
}

func (Recipe) TableName() string     { return "recipes" }
func (RecipeBase) TableName() string { return "recipes" }
