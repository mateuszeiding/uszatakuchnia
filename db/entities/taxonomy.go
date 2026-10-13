package entities

// TaxonomyAxis — oś klasyfikacji przepisów (Typ dania, Baza, Kuchnia, Diety, Praktyczne).
// Nowe osie dodaje się wierszem w tabeli, bez zmiany modelu.
type TaxonomyAxis struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Code       string `gorm:"type:varchar(32);not null;uniqueIndex:ux_taxonomy_axes_code"`
	Label      string `gorm:"type:varchar(64);not null"`
	Multi      bool   `gorm:"not null;default:false"`
	Required   bool   `gorm:"not null;default:false"`
	HasColor   bool   `gorm:"not null;default:false"`
	FilterMode string `gorm:"type:varchar(8);not null;default:'or'"` // 'and' | 'or'
	SortOrder  int    `gorm:"not null;default:0"`
}

func (TaxonomyAxis) TableName() string { return "taxonomy_axes" }

type TaxonomyValue struct {
	ID     uint         `gorm:"primaryKey;autoIncrement"`
	AxisID uint         `gorm:"not null;uniqueIndex:ux_taxonomy_values_axis_label,priority:1"`
	Axis   TaxonomyAxis `gorm:"foreignKey:AxisID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Label     string  `gorm:"type:varchar(64);not null;uniqueIndex:ux_taxonomy_values_axis_label,priority:2"`
	Color     *string `gorm:"type:varchar(16)"` // klucz tokenu koloru, np. "coral" → var(--c-coral)
	SortOrder int     `gorm:"not null;default:0"`
}

func (TaxonomyValue) TableName() string { return "taxonomy_values" }

type RecipeTaxonomy struct {
	RecipeID uint          `gorm:"not null;primaryKey"`
	ValueID  uint          `gorm:"not null;primaryKey;index:ix_recipe_taxonomy_value"`
	Value    TaxonomyValue `gorm:"foreignKey:ValueID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (RecipeTaxonomy) TableName() string { return "recipe_taxonomy" }
