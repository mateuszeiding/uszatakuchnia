package entities

type Ingredient struct {
	ID       uint        `gorm:"primaryKey;autoIncrement"`
	ParentID *uint       `gorm:"index:ix_ingredients_parent;uniqueIndex:ux_ingredients_parent_name"`
	Parent   *Ingredient `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Name   string         `gorm:"not null;uniqueIndex:ux_ingredients_parent_name"`
	TypeID uint           `gorm:"not null;index:ix_ingredients_type"`
	Type   IngredientType `gorm:"foreignKey:TypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	IsAllergen bool `gorm:"not null;default:false"`
}

func (Ingredient) TableName() string { return "ingredients" }
