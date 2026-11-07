package entities

type Ingredient struct {
	ID       uint        `gorm:"primaryKey;autoIncrement"`
	ParentID *uint       `gorm:"index:ix_ingredients_parent;uniqueIndex:ux_ingredients_parent_name"`
	Parent   *Ingredient `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Name   string         `gorm:"not null;uniqueIndex:ux_ingredients_parent_name"`
	TypeID uint           `gorm:"not null;index:ix_ingredients_type"`
	Type   IngredientType `gorm:"foreignKey:TypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Makro na 100 g, orientacyjnie (marki się różnią). *Min nil = brak własnych danych (dziedziczy z podskładników),
	// *Max nil = jedna wartość.
	KcalMin    *float64 `gorm:"type:numeric"`
	KcalMax    *float64 `gorm:"type:numeric"`
	ProteinMin *float64 `gorm:"type:numeric"`
	ProteinMax *float64 `gorm:"type:numeric"`
	FatMin     *float64 `gorm:"type:numeric"`
	FatMax     *float64 `gorm:"type:numeric"`
	CarbsMin   *float64 `gorm:"type:numeric"`
	CarbsMax   *float64 `gorm:"type:numeric"`
	FiberMin   *float64 `gorm:"type:numeric"`
	FiberMax   *float64 `gorm:"type:numeric"`

	Allergens   []IngredientAllergen       `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Conversions []IngredientUnitConversion `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Image Image `gorm:"polymorphic:Ref;polymorphicValue:ingredient"`
}

func (Ingredient) TableName() string { return "ingredients" }

// macroFields — wskaźniki na pary min/max po kluczu z NutrientKeys.
func (i *Ingredient) macroFields(key string) (min, max **float64) {
	switch key {
	case "kcal":
		return &i.KcalMin, &i.KcalMax
	case "protein":
		return &i.ProteinMin, &i.ProteinMax
	case "fat":
		return &i.FatMin, &i.FatMax
	case "carbs":
		return &i.CarbsMin, &i.CarbsMax
	case "fiber":
		return &i.FiberMin, &i.FiberMax
	}
	return nil, nil
}

// MacroOf — własna wartość na 100 g; nil = brak danych.
func (i Ingredient) MacroOf(key string) *MacroRange {
	min, max := i.macroFields(key)
	if min == nil || *min == nil {
		return nil
	}
	return &MacroRange{Min: **min, Max: *max}
}

// SetMacro — ustawia (r != nil) albo czyści (r == nil) parę min/max.
func (i *Ingredient) SetMacro(key string, r *MacroRange) {
	min, max := i.macroFields(key)
	if min == nil {
		return
	}
	if r == nil {
		*min, *max = nil, nil
		return
	}
	v := r.Min
	*min, *max = &v, r.Max
}
