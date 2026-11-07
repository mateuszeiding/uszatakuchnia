package entities

// IngredientUnitConversion — ile gramów waży jednostka danego składnika (np. jajko: 1 sztuka ≈ 50–65 g).
// g / kg / ml / l przeliczane w kodzie (services.GramsPerUnit), bez wierszy.
type IngredientUnitConversion struct {
	ID           uint     `gorm:"primaryKey;autoIncrement"`
	IngredientID uint     `gorm:"not null;uniqueIndex:ux_ing_unit_conv"`
	UnitID       uint     `gorm:"not null;uniqueIndex:ux_ing_unit_conv"`
	Unit         Unit     `gorm:"foreignKey:UnitID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	GramsMin     float64  `gorm:"type:numeric;not null"`
	GramsMax     *float64 `gorm:"type:numeric"`
}

func (IngredientUnitConversion) TableName() string { return "ingredient_unit_conversions" }
