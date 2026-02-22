package entities

import dtos "uszatakuchnia/dtos"

type Ingredient struct {
	ID       uint        `gorm:"primaryKey;autoIncrement"`
	ParentID *uint       `gorm:"index:ix_ingredients_parent;uniqueIndex:ux_ingredients_parent_name"`
	Parent   *Ingredient `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Name   string         `gorm:"not null;uniqueIndex:ux_ingredients_parent_name"`
	TypeID uint           `gorm:"not null;index:ix_ingredients_type"`
	Type   IngredientType `gorm:"foreignKey:TypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	IsAllergen bool `gorm:"not null;default:false"`

	Image Image `gorm:"polymorphic:Ref;polymorphicValue:ingredient"`
}

func (Ingredient) TableName() string { return "ingredients" }

func (i Ingredient) ToDto() dtos.IngredientDto {
	ingredientType := dtos.TypeDto{
		Name: i.Type.Name,
		Code: i.Type.Code,
	}

	return dtos.IngredientDto{
		ID:         i.ID,
		Name:       i.Name,
		ParentID:   i.ParentID,
		IsAllergen: i.IsAllergen,
		Type:       ingredientType,
		Image:      i.Image.ToDto(),
	}
}

type IngredientList []Ingredient

func (l IngredientList) ListToDto() []dtos.IngredientDto {
	out := make([]dtos.IngredientDto, 0, len(l))
	for _, i := range l {
		out = append(out, i.ToDto())
	}
	return out
}
