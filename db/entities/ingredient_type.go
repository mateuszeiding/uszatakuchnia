package entities

type IngredientType struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Code string `gorm:"type:varchar(64);not null;uniqueIndex:ux_ingredient_type_code"`
	Name string `gorm:"type:varchar(128);not null"`
}

func (IngredientType) TableName() string { return "ingredient_types" }
