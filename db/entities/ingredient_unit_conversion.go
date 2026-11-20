package entities

// g / kg are converted in code (services.GramsPerUnit), without rows. ml / l rows are optional
// overrides; without one, ml / l are derived from the other volume unit or fall back to water density.
type IngredientUnitConversion struct {
	ID           uint     `gorm:"primaryKey;autoIncrement"`
	IngredientID uint     `gorm:"not null;uniqueIndex:ux_ing_unit_conv"`
	UnitID       uint     `gorm:"not null;uniqueIndex:ux_ing_unit_conv"`
	Unit         Unit     `gorm:"foreignKey:UnitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	GramsMin     float64  `gorm:"type:numeric;not null"`
	GramsMax     *float64 `gorm:"type:numeric"`
}

func (IngredientUnitConversion) TableName() string { return "ingredient_unit_conversions" }
