package entities

type Aroma struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	IngredientID uint       `gorm:"not null;index:ix_aromas_ingredient;uniqueIndex:ux_aromas_ing_name"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Name   string    `gorm:"not null;uniqueIndex:ux_aromas_ing_name"`
	TypeID uint      `gorm:"not null;index:ix_aromas_type"`
	Type   AromaType `gorm:"foreignKey:TypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Intensity int `gorm:"not null;default:50"`
}

func (Aroma) TableName() string { return "aromas" }
