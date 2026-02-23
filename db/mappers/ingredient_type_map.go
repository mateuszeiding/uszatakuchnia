package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapIngredientTypeToDto(e entities.IngredientType) dtos.TypeDto {
	return dtos.TypeDto{
		Name: e.Name,
		Code: e.Code,
	}
}

func MapIngredientTypeArrayToDto(a []entities.IngredientType) []dtos.TypeDto {
	out := make([]dtos.TypeDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapIngredientTypeToDto(i))
	}
	return out
}
