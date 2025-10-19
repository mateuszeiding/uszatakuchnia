package entities

import (
	dtos "uszatakuchnia/dtos"
)

type Ingredient struct {
	ID       uint        `gorm:"primaryKey;autoIncrement"`
	ParentID *uint       `gorm:"index:ix_ingredients_parent;uniqueIndex:ux_ingredients_parent_name"`
	Parent   *Ingredient `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Name   string         `gorm:"not null;uniqueIndex:ux_ingredients_parent_name"`
	TypeID uint           `gorm:"not null;index:ix_ingredients_type"`
	Type   IngredientType `gorm:"foreignKey:TypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	IsAllergen bool `gorm:"not null;default:false"`

	Aromas           []Aroma
	IngredientTastes []IngredientTaste
}

func (Ingredient) TableName() string { return "ingredients" }

func (i Ingredient) ToDto() dtos.IngredientDto {
	aromas := make([]dtos.AromaDto, 0, len(i.Aromas))
	for _, a := range i.Aromas {
		aromas = append(aromas, dtos.AromaDto{
			Name:      a.Name,
			Intensity: a.Intensity,
			Code:      a.Type.Code,
		})
	}

	tastes := make([]dtos.TasteDto, 0, len(i.IngredientTastes))
	for _, t := range i.IngredientTastes {
		tastes = append(tastes, dtos.TasteDto{
			Name:      t.Taste.Type.Name,
			Intensity: t.Intensity,
			Code:      t.Taste.Type.Code,
		})
	}

	dto := dtos.IngredientDto{
		ID:         i.ID,
		Name:       i.Name,
		ParentID:   i.ParentID,
		IsAllergen: i.IsAllergen,
		Type:       i.Type.Code,
		Aromas:     aromas,
		Tastes:     tastes,
	}

	return dto
}

type IngredientList []Ingredient

func (l IngredientList) ListToDto() []dtos.IngredientDto {
	dtos := make([]dtos.IngredientDto, 0, len(l))
	for _, i := range l {
		dtos = append(dtos, i.ToDto())
	}

	return dtos
}
