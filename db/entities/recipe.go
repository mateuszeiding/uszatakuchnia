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

	Photo       *RecipePhoto       `gorm:"foreignKey:RecipeID;references:ID"`
	Steps       []RecipeStep       `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Ingredients []RecipeIngredient `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type RecipeBase struct {
	ID          uint         `gorm:"primaryKey;autoIncrement"`
	Name        string       `gorm:"type:varchar(255);not null;index:ix_recipes_name"`
	Tagline     *string      `gorm:"type:varchar(500)"`
	Category    *string      `gorm:"type:varchar(64)"`
	Region      *string      `gorm:"type:varchar(64)"`
	TimeMinutes *int
	Difficulty  *int
	Photo       *RecipePhoto `gorm:"foreignKey:RecipeID;references:ID"`
}

func (Recipe) TableName() string { return "recipes" }
