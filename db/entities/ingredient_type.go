package entities

import "uszatakuchnia/dtos"

type IngredientType struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Code string `gorm:"type:varchar(64);not null;uniqueIndex:ux_ingredient_type_code"`
	Name string `gorm:"type:varchar(128);not null"`
}

func (IngredientType) TableName() string { return "ingredient_types" }

func (i IngredientType) ToDto() dtos.TypeDto {
	return dtos.TypeDto{
		Name: i.Name,
		Code: i.Code,
	}
}

type IngredientTypeList []IngredientType

func (l IngredientTypeList) ListToDto() []dtos.TypeDto {
	out := make([]dtos.TypeDto, 0, len(l))
	for _, i := range l {
		out = append(out, i.ToDto())
	}
	return out
}
