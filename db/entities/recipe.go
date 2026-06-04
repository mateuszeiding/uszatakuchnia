package entities

type Recipe struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"type:varchar(255);not null;index:ix_recipes_name"`
	Servings    int     `gorm:"not null"`
	Description *string
	Tagline     *string `gorm:"type:varchar(500)"`

	Category       *string `gorm:"type:varchar(64)"`
	Region         *string `gorm:"type:varchar(64)"`
	TimeMinutes    *int
	Difficulty     *int
	KcalPerServing *int

	Status    string `gorm:"type:varchar(16);not null;default:'published'"`
	NeedsPrep bool   `gorm:"not null;default:false"`

	Photo       *RecipePhoto       `gorm:"foreignKey:RecipeID;references:ID"`
	Steps       []RecipeStep       `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Ingredients []RecipeIngredient `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tags        []RecipeTag        `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type RecipeBase struct {
	ID          uint         `gorm:"primaryKey;autoIncrement"`
	Name        string       `gorm:"type:varchar(255);not null;index:ix_recipes_name"`
	Tagline     *string      `gorm:"type:varchar(500)"`
	Category    *string      `gorm:"type:varchar(64)"`
	Region      *string      `gorm:"type:varchar(64)"`
	TimeMinutes *int
	Difficulty  *int
	Status      string       `gorm:"type:varchar(16);not null;default:'published'"`
	NeedsPrep   bool         `gorm:"not null;default:false"`
	Photo       *RecipePhoto `gorm:"foreignKey:RecipeID;references:ID"`
	Tags        []RecipeTag  `gorm:"foreignKey:RecipeID;references:ID"`
}

func (Recipe) TableName() string    { return "recipes" }
func (RecipeBase) TableName() string { return "recipes" }
